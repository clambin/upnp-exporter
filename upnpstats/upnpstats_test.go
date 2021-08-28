package upnpstats_test

import (
	"github.com/clambin/upnp-exporter/upnpstats"
	"github.com/stretchr/testify/assert"
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
