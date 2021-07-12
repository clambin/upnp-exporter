package main

import (
	"fmt"
	"github.com/clambin/upnp-exporter/internal/upnpstats"
	"github.com/huin/goupnp/dcps/internetgateway1"
	log "github.com/sirupsen/logrus"
)

func main() {
	routers, err := upnpstats.DiscoverRouters()

	if err != nil {
		log.WithError(err).Fatal("unable to discover router(s)")
	}

	for _, router := range routers {
		var clients []*internetgateway1.WANCommonInterfaceConfig1

		clients, err = internetgateway1.NewWANCommonInterfaceConfig1ClientsByURL(&router.URLBase)

		if err != nil {
			log.WithError(err).Warning("failed to create clients from router URL")
		}

		for _, client := range clients {
			fmt.Println(client.GetCommonLinkProperties())
		}
	}
}
