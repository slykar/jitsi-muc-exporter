package main

import (
	"crypto/tls"
	"gopkg.in/alecthomas/kingpin.v2"
	"gosrc.io/xmpp"
)

type Config struct {
	Host        string
	Port        int
	AuthDomain  string
	Username    string
	Password    string
	Brewery     string
	MucNickname string
	SkipTls     bool
}

func (c *Config) GetXMPPConfig() *xmpp.Config {
	return &xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address:        c.Host,
			Domain:         c.AuthDomain,
			ConnectTimeout: 10,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: c.SkipTls,
			},
		},
		Jid:            c.Username + "@" + c.AuthDomain,
		Credential:     xmpp.Password(c.Password),
		ConnectTimeout: 10,
	}
}

func Configure(app *kingpin.Application) *Config {
	c := &Config{}

	// You can provide hint options statically
	app.Flag("host", "Host address of the Prosody XMPP server.").Short('h').
		Required().
		StringVar(&c.Host)

	app.Flag("port", "Prosody XMPP server port.").
		Default("5222").
		IntVar(&c.Port)

	app.Flag("domain", "XMPP auth domain. Should be the same as for your JVB.").
		Required().
		StringVar(&c.AuthDomain)

	app.Flag("brewery", "MUC JID of the JVB Brewery room. E.g. jvbbrewery@internal-muc.meet.jitsi.").
		Required().
		StringVar(&c.Brewery)

	app.Flag("username", "XMPP username. Can be shared with your JVB.").Short('u').
		Default("jvb").
		StringVar(&c.Username)

	app.Flag("password", "XMPP password. Can be shared with your JVB.").
		Required().
		StringVar(&c.Password)

	app.Flag("nickname", "MUC nickname for the exporter.").
		Default("jitsi-muc-exporter").
		StringVar(&c.MucNickname)

	app.Flag("insecure", "Skip TLS verification when connecting to Prosody.").
		BoolVar(&c.SkipTls)

	return c
}
