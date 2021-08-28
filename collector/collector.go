package collector

import (
	"github.com/clambin/upnp-exporter/upnpstats"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Collector struct {
	Scanner         upnpstats.Scanner
	packetsSent     *prometheus.Desc
	packetsReceived *prometheus.Desc
	bytesSent       *prometheus.Desc
	bytesReceived   *prometheus.Desc
}

func New(scanner upnpstats.Scanner) *Collector {
	return &Collector{
		Scanner: scanner,
		packetsSent: prometheus.NewDesc(
			prometheus.BuildFQName("upnp", "stats", "sent_packets"),
			"Total number of packets sent by the router",
			[]string{"router"},
			nil),
		packetsReceived: prometheus.NewDesc(
			prometheus.BuildFQName("upnp", "stats", "received_packets"),
			"Total number of packets received by the router",
			[]string{"router"},
			nil),
		bytesSent: prometheus.NewDesc(
			prometheus.BuildFQName("upnp", "stats", "sent_bytes"),
			"Total number of bytes sent by the router",
			[]string{"router"},
			nil),
		bytesReceived: prometheus.NewDesc(
			prometheus.BuildFQName("upnp", "stats", "received_bytes"),
			"Total number of bytes received by the router",
			[]string{"router"},
			nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.packetsSent
	ch <- c.packetsReceived
	ch <- c.bytesSent
	ch <- c.bytesReceived
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.Scanner.ReportNetworkStats()

	if err != nil {
		log.WithError(err).Error("failed to scan router(s)")
		return
	}

	for _, routerStats := range stats {
		router := routerStats.RouterURL
		ch <- prometheus.MustNewConstMetric(c.packetsSent, prometheus.CounterValue, float64(routerStats.PacketsSent), router)
		ch <- prometheus.MustNewConstMetric(c.packetsReceived, prometheus.CounterValue, float64(routerStats.PacketsReceived), router)
		ch <- prometheus.MustNewConstMetric(c.bytesSent, prometheus.CounterValue, float64(routerStats.BytesSent), router)
		ch <- prometheus.MustNewConstMetric(c.bytesReceived, prometheus.CounterValue, float64(routerStats.BytesReceived), router)
	}
}
