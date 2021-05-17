package upnpstats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	routerPacketsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "upnp_stats_sent_packets",
		Help: "Total number of packets sent by the router",
	}, []string{"router"})
	routerPacketsReceived = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "upnp_stats_received_packets",
		Help: "Total number of packets received by the router",
	}, []string{"router"})
	routerBytesSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "upnp_stats_sent_bytes",
		Help: "Total number of bytes sent by the router",
	}, []string{"router"})
	routerBytesReceived = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "upnp_stats_received_bytes",
		Help: "Total number of bytes received by the router",
	}, []string{"router"})
)
