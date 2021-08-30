package upnpstats_test

import (
	"fmt"
	"github.com/clambin/upnp-exporter/upnpstats"
	"github.com/clambin/upnp-exporter/upnpstats/mocks"
	"github.com/huin/goupnp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestRouterScanner_ReportNetworkStats(t *testing.T) {
	scanner, err := upnpstats.New(nil)
	if err != nil {
		t.Log("no routers found. Not running test")
		return
	}

	stats, err := scanner.ReportNetworkStats()

	for _, routerStats := range stats {
		t.Logf("testing router %s", routerStats.RouterURL)
		assert.NotZero(t, routerStats.PacketsSent)
		assert.NotZero(t, routerStats.PacketsReceived)
		assert.NotZero(t, routerStats.BytesSent)
		assert.NotZero(t, routerStats.BytesReceived)
	}
}

func TestRouterScanner_Routers(t *testing.T) {
	myURL, _ := url.Parse("https://127.0.0.1")
	scanner, err := upnpstats.New(myURL)
	require.NoError(t, err)

	routers := scanner.Routers()
	require.Len(t, routers, 1)
	assert.Equal(t, "https://127.0.0.1", routers[0])
}

func TestRouterScanner_ReportNetworkStats_Mocked(t *testing.T) {
	discoverer := &mocks.Discoverer{}
	scanner := &upnpstats.RouterScanner{
		Discoverer: discoverer,
	}

	const myURL = "https://127.0.0.1:5000/foo"
	routerURL, _ := url.Parse(myURL)
	discoverer.
		On("DiscoverDevices", "urn:schemas-upnp-org:device:InternetGatewayDevice:1").
		Return([]goupnp.MaybeRootDevice{
			{Location: routerURL, Root: &goupnp.RootDevice{URLBase: *routerURL}},
		}, nil).
		Once()

	err := scanner.Discover()
	require.NoError(t, err)

	discoverer.
		On("GetNetworkStats", mock.AnythingOfType("*url.URL")).
		Return(upnpstats.Stats{
			RouterURL:       myURL,
			PacketsSent:     10,
			PacketsReceived: 20,
			BytesSent:       30,
			BytesReceived:   40,
		}, nil).Once()

	stats, err := scanner.ReportNetworkStats()
	require.NoError(t, err)
	assert.Len(t, stats, 1)

	for _, routerStats := range stats {
		t.Logf("testing router %s", routerStats.RouterURL)
		assert.Equal(t, uint32(10), routerStats.PacketsSent)
		assert.Equal(t, uint32(20), routerStats.PacketsReceived)
		assert.Equal(t, uint64(30), routerStats.BytesSent)
		assert.Equal(t, uint64(40), routerStats.BytesReceived)
	}
}

func TestRouterScanner_ReportNetworkStats_Mocked_Failures(t *testing.T) {
	discoverer := &mocks.Discoverer{}
	scanner := &upnpstats.RouterScanner{
		Discoverer: discoverer,
	}

	discoverer.
		On("DiscoverDevices", "urn:schemas-upnp-org:device:InternetGatewayDevice:1").
		Return(nil, fmt.Errorf("failed to discover devices")).
		Once()

	err := scanner.Discover()
	require.Error(t, err)

	const myURL = "https://127.0.0.1:5000/foo"
	routerURL, _ := url.Parse(myURL)
	discoverer.
		On("DiscoverDevices", "urn:schemas-upnp-org:device:InternetGatewayDevice:1").
		Return([]goupnp.MaybeRootDevice{
			{Location: routerURL, Root: &goupnp.RootDevice{URLBase: *routerURL}},
		}, nil).
		Once()

	err = scanner.Discover()
	require.NoError(t, err)

	discoverer.
		On("GetNetworkStats", mock.AnythingOfType("*url.URL")).
		Return(upnpstats.Stats{}, fmt.Errorf("unable to get statistics")).Once()

	_, err = scanner.ReportNetworkStats()
	require.Error(t, err)
}
