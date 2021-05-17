package upnpstats

import (
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
	log "github.com/sirupsen/logrus"
	"net/url"
	"time"
)

func DiscoverRouterURLs() (urls []*url.URL, err error) {
	var devices []goupnp.MaybeRootDevice

	devices, err = goupnp.DiscoverDevices("urn:schemas-upnp-org:device:InternetGatewayDevice:1")

	if err != nil {
		log.WithError(err).Warning("unable to discover router URLS")
		return
	}

	for _, device := range devices {
		if device.Err == nil {
			log.WithField("device", device.Location.String()).Debug("router found")
			urls = append(urls, device.Location)
		}
	}

	return
}

func ReportNetworkStats(routerURL *url.URL) {
	clients, err := internetgateway1.NewWANCommonInterfaceConfig1ClientsByURL(routerURL)

	if err != nil {
		log.WithError(err).Error("unable to get clients")
		return
	}

	for _, client := range clients {
		var packets uint32

		packets, err = client.GetTotalPacketsReceived()

		if err == nil {
			routerPacketsReceived.WithLabelValues(client.RootDevice.URLBase.Host).Set(float64(packets))
			log.Debugf("packets received: %d", packets)
		} else {
			log.WithError(err).Warning("unable to get number of packets received")
		}

		packets, err = client.GetTotalPacketsSent()

		if err == nil {
			routerPacketsSent.WithLabelValues(client.RootDevice.URLBase.Host).Set(float64(packets))
			log.Debugf("packets sent: %d", packets)
		} else {
			log.WithError(err).Warning("unable to get number of packets received")
		}

		var bytes uint64

		bytes, err = client.GetTotalBytesReceived()

		if err == nil {
			routerBytesReceived.WithLabelValues(client.RootDevice.URLBase.Host).Set(float64(bytes))
			log.Debugf("bytes received: %d", bytes)
		} else {
			log.WithError(err).Warning("unable to get number of bytes received")
		}

		bytes, err = client.GetTotalBytesSent()

		if err == nil {
			routerBytesSent.WithLabelValues(client.RootDevice.URLBase.Host).Set(float64(bytes))
			log.Debugf("bytes sent: %d", bytes)
		} else {
			log.WithError(err).Warning("unable to get number of bytes received")
		}
	}
}

func Run(router *url.URL, interval time.Duration) {
	var routers []*url.URL
	var err error

	if router != nil {
		routers = append(routers, router)
	} else {
		if routers, err = DiscoverRouterURLs(); err != nil {
			log.WithField("err", err).Fatal("unable to discover routers. exiting")
		}
	}

	for {
		for _, router := range routers {
			ReportNetworkStats(router)
		}
		time.Sleep(interval)
	}
}
