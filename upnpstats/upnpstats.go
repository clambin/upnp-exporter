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
	Discoverer Discoverer
	routers    []url.URL
}

// New creates a new RouterScanner.  If router is provided, only that URL will be scanned. Otherwise, New will
// scan the network for compatible routers.
func New(router *url.URL) (scanner *RouterScanner, err error) {
	scanner = &RouterScanner{
		Discoverer: &UPNPDiscoverer{},
	}

	if router == nil {
		err = scanner.Discover()
	} else {
		scanner.routers = []url.URL{*router}
	}
	return
}

// Discover finds all upnp compatible routers
func (scanner *RouterScanner) Discover() (err error) {
	scanner.routers, err = scanner.discoverRouters()
	return
}

// discoverRouters attempts to discover all upnp-compatible routers
func (scanner *RouterScanner) discoverRouters() (routers []url.URL, err error) {
	var devices []goupnp.MaybeRootDevice

	devices, err = scanner.Discoverer.DiscoverDevices("urn:schemas-upnp-org:device:InternetGatewayDevice:1")

	if err != nil {
		log.WithError(err).Warning("unable to discover router URLS")
		return
	}

	for _, device := range devices {
		if device.Err == nil {
			log.WithField("device", device.Location.String()).Debug("router found")
			routers = append(routers, device.Root.URLBase)
		}
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
		var routerStats Stats
		routerStats, err = scanner.Discoverer.GetNetworkStats(&router)

		if err == nil {
			stats = append(stats, routerStats)
		} else {
			log.WithError(err).Warningf("failed to retrieve stats for %s", router.String())
		}
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

// Discoverer interface scans upnp-enabled routers
//go:generate mockery --name Discoverer
type Discoverer interface {
	DiscoverDevices(target string) (devices []goupnp.MaybeRootDevice, err error)
	GetNetworkStats(router *url.URL) (stats Stats, err error)
}

type UPNPDiscoverer struct {
}

func (discoverer *UPNPDiscoverer) DiscoverDevices(target string) (devices []goupnp.MaybeRootDevice, err error) {
	return goupnp.DiscoverDevices(target)
}

func (discoverer *UPNPDiscoverer) GetNetworkStats(routerURL *url.URL) (stats Stats, err error) {
	clients, err := internetgateway1.NewWANCommonInterfaceConfig1ClientsByURL(routerURL)

	if err != nil {
		log.WithError(err).Error("unable to get clients")
		return
	}

	stats.RouterURL = routerURL.String()
	for _, client := range clients {
		stats.PacketsReceived, err = client.GetTotalPacketsReceived()
		if err != nil {
			continue
		}

		stats.PacketsSent, err = client.GetTotalPacketsSent()
		if err != nil {
			continue
		}

		stats.BytesReceived, err = client.GetTotalBytesReceived()
		if err != nil {
			continue
		}

		stats.BytesSent, err = client.GetTotalBytesSent()
		if err != nil {
			continue
		}
	}

	return
}
