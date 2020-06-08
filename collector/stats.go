package collector

import "errors"

// Possible type of the stat that JVB exposes
type StatType int

const (
	_ StatType = iota
	StatCounter
	StatGauge
	StatHistogram
	StatTag
)

// Describes the stat that JVB exposes, so we can convert it to a metric understood by Prometheus.
type StatDescriptor struct {
	Type StatType
	Help string
}

// A mapping of supported stats and their respective types that we can map to something
// that Prometheus can understand. An optional help text is also provided.
//
// I've only added stats that I'm interested in. Feel free to extend this list.
// I'm also missing a lot of the help texts. Some of them can be found under the link below,
// but they seem to be outdated at the moment.
//
// https://github.com/jitsi/jitsi-videobridge/blob/master/doc/statistics.md
//
// TODO: Add all possible stats
// TODO: Add help text
var jvbStatsDescriptors = map[string]StatDescriptor{
	"version": {StatTag, ""},
	"threads": {StatGauge, ""},

	"p2p_conferences":    {StatGauge, ""},
	"conferences":        {StatGauge, ""},
	"participants":       {StatGauge, ""},
	"videostreams":       {StatGauge, ""},
	"videochannels":      {StatGauge, ""},
	"largest_conference": {StatGauge, ""},

	"endpoints_sending_video": {StatGauge, ""},
	"endpoints_sending_audio": {StatGauge, ""},

	"bit_rate_download": {StatGauge, ""},
	"bit_rate_upload":   {StatGauge, ""},

	"conference_sizes":             {StatHistogram, "The distribution of conference sizes hosted on the bridge."},
	"conferences_by_video_senders": {StatHistogram, ""},
	"conferences_by_audio_senders": {StatHistogram, ""},

	"total_participants":       {StatCounter, ""},
	"total_conference_seconds": {StatCounter, ""},
}

// Get a descriptor by JVB stat name.
// This descriptor is helpful to identify the related Prometheus metric type.
func GetStatDescriptor(statName string) (StatDescriptor, error) {
	if value, ok := jvbStatsDescriptors[statName]; ok {
		return value, nil
	} else {
		return value, errors.New("descriptor for stat not found")
	}
}
