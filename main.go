package main

import (
	"encoding/xml"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slykar/jitsi-muc-exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"net/http"
	"os"
)

// Creates a new Presence that needs to be sent after connecting to the server in order
// to start receiving JVB Presence packets from the JVB brewery room.
func NewJvbBreweryPresence(room, nickname string) stanza.Presence {
	presence := stanza.NewPresence(stanza.Attrs{
		To: room + "/" + nickname,
	})

	presence.Extensions = append(presence.Extensions, stanza.MucPresence{})

	return presence
}

var (
	jvbCollector = collector.NewJvbMucCollector()
)

func init() {
	// Register extension for JVB Stats element
	stanza.TypeRegistry.MapExtension(
		stanza.PKTPresence,
		xml.Name{Space: "http://jitsi.org/protocol/colibri", Local: "stats"},
		collector.JvbStats{},
	)

	prometheus.MustRegister(jvbCollector)
}

func main() {
	app := kingpin.New("jitsi-muc-exporter", "Export all JVB stats by joining JVB Brewery MUC room.")
	config := Configure(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// Router is required to assign handlers for different events
	router := xmpp.NewRouter()

	// The only handler we need to implement is the "presence" handler
	router.HandleFunc("presence", func(s xmpp.Sender, p stanza.Packet) {
		presence, ok := p.(stanza.Presence)

		if !ok {
			log.Println("Could not cast presence packet as stanza.Presence.")
		}

		stats := collector.JvbStats{}

		if presence.Get(&stats) {
			log.Printf("Presence packet received from %s\n", presence.From)
			jvbCollector.UpdateWithStats(collector.IdentifyJvbFromPresence(presence), &stats)
		}
	})

	// Create client instance - this will not connect yet.
	// A StreamManager is created below to manage the connection.
	client, err := xmpp.NewClient(config.GetXMPPConfig(), router, func(err error) {
		log.Fatalln("Could not create Client instance")
	})

	if err != nil {
		log.Fatalln("Could not connect to the XMPP server")
	}

	// If you pass the client to a connection manager, it will handle the reconnect policy
	// for you automatically.
	streamManager := xmpp.NewStreamManager(client, func(c xmpp.Sender) {
		// After connecting, we need to join the "brewery" room,
		// where JVBs communicate their metrics as presence packets
		client := c.(*xmpp.Client)

		// Say hello to everyone in order to receive presence updates.
		// TODO: Configurable MUC nickname for the exporter
		err := client.Send(NewJvbBreweryPresence(config.Brewery, config.MucNickname))

		if err != nil {
			log.Fatalln("Could not send presence packet to JVB brewery room.")
		}

	})

	// Run XMPP monitoring in a coroutine, we will block the process with Prometheus HTTP server
	go func() { log.Fatal(streamManager.Run()) }()

	// Add Prom HTTP metrics handler
	log.Println("Starting metrics server on port :2112")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
