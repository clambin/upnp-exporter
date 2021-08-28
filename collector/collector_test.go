package collector_test

import (
	"fmt"
	"github.com/clambin/gotools/metrics"
	"github.com/clambin/upnp-exporter/collector"
	"github.com/clambin/upnp-exporter/upnpstats"
	"github.com/clambin/upnp-exporter/upnpstats/mocks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollector_Describe(t *testing.T) {
	c := collector.New(nil)

	ch := make(chan *prometheus.Desc)
	go c.Describe(ch)

	for _, metricName := range []string{
		"upnp_stats_sent_packets",
		"upnp_stats_received_packets",
		"upnp_stats_sent_bytes",
		"upnp_stats_received_bytes",
	} {
		metric := <-ch
		assert.Contains(t, metric.String(), "\""+metricName+"\"")
	}
}

func TestCollector_Collect(t *testing.T) {
	const router = "https://127.0.0.1"
	scanner := &mocks.Scanner{}
	c := collector.New(scanner)

	scanner.
		On("ReportNetworkStats").
		Return([]upnpstats.Stats{{
			RouterURL:       router,
			PacketsSent:     1,
			PacketsReceived: 2,
			BytesSent:       10,
			BytesReceived:   20,
		}}, nil).
		Once()

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	testCases := []struct {
		metricName  string
		metricValue float64
	}{
		{metricName: "upnp_stats_sent_packets", metricValue: 1.0},
		{metricName: "upnp_stats_received_packets", metricValue: 2.0},
		{metricName: "upnp_stats_sent_bytes", metricValue: 10.0},
		{metricName: "upnp_stats_received_bytes", metricValue: 20.0},
	}

	for _, testCase := range testCases {
		m := <-ch
		assert.Equal(t, testCase.metricName, metrics.MetricName(m))
		assert.Equal(t, testCase.metricValue, metrics.MetricValue(m).GetCounter().GetValue())
		assert.Equal(t, router, metrics.MetricLabel(m, "router"))
	}
}

func TestCollector_Collect_Fail(t *testing.T) {
	const router = "https://127.0.0.1"
	scanner := &mocks.Scanner{}
	c := collector.New(scanner)

	scanner.
		On("ReportNetworkStats").
		Return([]upnpstats.Stats{}, fmt.Errorf("unable to scan %s", router)).
		Once()

	ch := make(chan prometheus.Metric)
	go c.Collect(ch)

	assert.Never(t, func() bool {
		<-ch
		return true
	}, 100*time.Millisecond, 10*time.Millisecond)
}
