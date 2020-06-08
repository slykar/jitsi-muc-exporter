package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

func parseConstMetric(stat JvbStat, desc *prometheus.Desc, valType prometheus.ValueType) (prometheus.Metric, error) {
	if value, err := strconv.ParseFloat(stat.Value, 32); err == nil {
		return prometheus.NewConstMetric(desc, valType, value)
	} else {
		return nil, err
	}
}

func ParseGauge(stat JvbStat, desc *prometheus.Desc) (prometheus.Metric, error) {
	return parseConstMetric(stat, desc, prometheus.GaugeValue)
}

func ParseCounter(stat JvbStat, desc *prometheus.Desc) (prometheus.Metric, error) {
	return parseConstMetric(stat, desc, prometheus.CounterValue)
}

func ParseHistogram(stat JvbStat, desc *prometheus.Desc) (prometheus.Metric, error) {

	// try to parse JSON array value that represents our bucketed values
	bucketedValues, err := ParseArray(stat.Value)

	if err != nil {
		return nil, err
	}

	var (
		// we will have as many buckets as we got from the stats
		buckets = make(map[float64]uint64, len(bucketedValues))
		count   = uint64(0)
		sum     = float64(0)
	)

	for bucketSize, bucketCount := range bucketedValues {
		// Prometheus histograms are cumulative - increase count and then set bucket value
		count += bucketCount
		buckets[float64(bucketSize)] = count
		sum += float64(bucketSize) * float64(bucketCount)
	}

	return prometheus.NewConstHistogram(desc, count, sum, buckets)
}
