package upnpstats

import (
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
	log "github.com/sirupsen/logrus"
	"net/url"
)

// Scanner interface scans upnp-enabled routers for network statistics
//go:generate mockery --name Scanner
type Scanner interface {
	ReportNetworkStats() (stats []Stats, err error)
	Routers() (routers []string)
}

// RouterScanner scans upnp-enabled routers for network statistics
type RouterScanner struct {
	routers []url.URL
}

// New creates a new RouterScanner.  If router is provided, only that URL will be scanned. Otherwise, New will
// scan the network for compatible routers.
func New(router *url.URL) (scanner *RouterScanner, err error) {
	scanner = &RouterScanner{}
	if router == nil {
		err = scanner.discoverRouters()
	} else {
		scanner.routers = []url.URL{*router}
	}
	return
}

// Stats contains the statistics scanned from a router
type Stats struct {
	RouterURL       string
	PacketsSent     uint32
	PacketsReceived uint32
	BytesSent       uint64
	BytesReceived   uint64
}

// ReportNetworkStats scans all routers for updated network statistics
func (scanner *RouterScanner) ReportNetworkStats() (stats []Stats, err error) {
	for _, router := range scanner.routers {
		routerStats := Stats{RouterURL: router.String()}
		routerStats.PacketsSent, routerStats.PacketsReceived, routerStats.BytesSent, routerStats.BytesReceived, err = reportNetworkStats(&router)

		if err != nil {
			log.WithError(err).Warningf("failed to retrieve stats for %s", router.String())
			continue
		}
		stats = append(stats, routerStats)
	}
	return
}

// Routers returns the list of routers
func (scanner *RouterScanner) Routers() (routers []string) {
	for _, r := range scanner.routers {
		routers = append(routers, r.String())
	}
	return
}

// discoverRouters attempts to discover all upnp-compatible routers
func (scanner *RouterScanner) discoverRouters() (err error) {
	var devices []goupnp.MaybeRootDevice

	devices, err = goupnp.DiscoverDevices("urn:schemas-upnp-org:device:InternetGatewayDevice:1")

	if err != nil {
		log.WithError(err).Warning("unable to discover router URLS")
		return
	}

	for _, device := range devices {
		if device.Err == nil {
			log.WithField("device", device.Location.String()).Debug("router found")
			scanner.routers = append(scanner.routers, device.Root.URLBase)
		}
	}

	return
}

func reportNetworkStats(routerURL *url.URL) (packetsSent, packetsReceived uint32, bytesSent, bytesReceived uint64, err error) {
	clients, err := internetgateway1.NewWANCommonInterfaceConfig1ClientsByURL(routerURL)

	if err != nil {
		log.WithError(err).Error("unable to get clients")
		return
	}

	for _, client := range clients {
		packetsReceived, err = client.GetTotalPacketsReceived()
		if err != nil {
			continue
		}

		packetsSent, err = client.GetTotalPacketsSent()
		if err != nil {
			continue
		}

		bytesReceived, err = client.GetTotalBytesReceived()
		if err != nil {
			continue
		}

		bytesSent, err = client.GetTotalBytesSent()
		if err != nil {
			continue
		}
	}

	return
}
