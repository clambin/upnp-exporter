package main

import (
	"fmt"
	"github.com/clambin/upnp-exporter/upnpstats"
	log "github.com/sirupsen/logrus"
)

func main() {
	scanner, err := upnpstats.New(nil)

	if err != nil {
		log.WithError(err).Fatal("unable to discover router(s)")
	}

	stats, err := scanner.ReportNetworkStats()

	for _, routerStats := range stats {
		fmt.Println(routerStats)
	}
}
