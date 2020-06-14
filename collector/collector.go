package collector

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"gosrc.io/xmpp/stanza"
	"log"
	"strings"
)

type StatsByName = map[string]JvbStat

type JvbMucCollector struct {
	namespace string
	// Keep stats for each JVB - MUC stats get reported for all connected JVBs
	// TODO: Should we keep a map of stats instead of a list?
	//  How does this affect GC if we just replace the pointer to JvbStats all the time?
	statsByJvb map[JvbIdentity]*JvbStats
}

type JvbStat struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type JvbStats struct {
	Stats []JvbStat `xml:"stat"`
}

// A string representing the "name" of the JVB.
// Right now it is a MUC nickname found in the presence packet.
type JvbIdentity string

func ParseArray(s string) ([]uint64, error) {
	var numbers []uint64
	return numbers, json.Unmarshal([]byte(s), &numbers)
}

func NewJvbMucCollector(namespace string) *JvbMucCollector {
	return &JvbMucCollector{
		namespace:  namespace,
		statsByJvb: make(map[JvbIdentity]*JvbStats),
	}
}

func NewPromDescForStat(namespace string, jvbId JvbIdentity, stat JvbStat, descriptor *StatDescriptor) *prometheus.Desc {
	constLabels := prometheus.Labels{"jvb": string(jvbId)}
	fqName := namespace + "_" + stat.Name
	return prometheus.NewDesc(fqName, descriptor.Help, []string{}, constLabels)
}

// This implements Prometheus interface for collectors
func (c *JvbMucCollector) Describe(ch chan<- *prometheus.Desc) {
	for jvbId, stats := range c.statsByJvb {
		for _, stat := range stats.Stats {
			if descriptor, err := GetStatDescriptor(stat.Name); err == nil {
				ch <- NewPromDescForStat(c.namespace, jvbId, stat, &descriptor)
			}
		}
	}
}

func (c *JvbMucCollector) Collect(ch chan<- prometheus.Metric) {

	for jvbId, stats := range c.statsByJvb {
		for _, stat := range stats.Stats {
			// Try to get a descriptor of the stat by it's name.
			// We need to know what type of value we are dealing with.
			descriptor, err := GetStatDescriptor(stat.Name)

			if err != nil {
				// not a known stat - continue to the next one
				continue
			}

			var (
				metric prometheus.Metric
				desc   = NewPromDescForStat(c.namespace, jvbId, stat, &descriptor)
			)

			switch descriptor.Type {
			case StatGauge:
				metric, err = ParseGauge(stat, desc)
			case StatCounter:
				metric, err = ParseCounter(stat, desc)
			case StatHistogram:
				metric, err = ParseHistogram(stat, desc)
			default:
				// type support not implemented - continue to the next stat
				continue
			}

			if metric != nil {
				ch <- metric
			} else {
				log.Printf("An error occurred when trying to get metric for %s\n%s\n", stat.Name, err)
			}
		}
	}

}

// Basically gets the MUC nickname of a JVB.
// It's very naive to return the element at idx 1, but the nickname should be there.
// TODO: Less naive approach? Only if there is a problem with current approach
// 	- I guess MUC nickname will always be there and not contain any "/" chars.
func IdentifyJvbFromPresence(p stanza.Presence) JvbIdentity {
	return JvbIdentity(strings.Split(p.From, "/")[1])
}

// Updates stats for a given JVB
func (c *JvbMucCollector) UpdateWithStats(jvbId JvbIdentity, stats *JvbStats) {

	log.Printf("Updating Metrics with Stats from %s", jvbId)

	// Simply add/replace stats for given JVB
	// TODO: Might be a good idea to at least filter the list of stats here,
	//  as we need to have it filtered for Collect and Describe anyway
	c.statsByJvb[jvbId] = stats
}
