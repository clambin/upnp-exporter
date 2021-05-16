package upnpstats

import (
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
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
)

func DiscoverRouters() (routers []*goupnp.RootDevice, err error) {
	var devices []goupnp.MaybeRootDevice

	devices, err = goupnp.DiscoverDevices("urn:schemas-upnp-org:device:InternetGatewayDevice:1")

	if err == nil {
		for _, device := range devices {
			log.WithField("device", device.Location.String()).Debug("router found")
			routers = append(routers, device.Root)
		}
	} else {
		log.WithError(err).Warning("unable to discover routers")
	}

	return
}

func ReportNetworkStats(router *goupnp.RootDevice) {
	clients, err := internetgateway1.NewWANCommonInterfaceConfig1ClientsFromRootDevice(router, nil)

	if err == nil {
		for _, client := range clients {
			var packets uint32

			packets, err = client.GetTotalPacketsReceived()

			if err == nil {
				routerPacketsReceived.WithLabelValues(client.RootDevice.URLBase.Host).Set(float64(packets))
			} else {
				log.WithError(err).Warning("unable to get number of packets received")
			}

			packets, err = client.GetTotalPacketsSent()

			if err == nil {
				routerPacketsSent.WithLabelValues(client.RootDevice.URLBase.Host).Set(float64(packets))
			} else {
				log.WithError(err).Warning("unable to get number of packets received")
			}
		}
	}
}
