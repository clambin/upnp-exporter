package main

import (
	"fmt"
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/ssdp"
	log "github.com/sirupsen/logrus"
)

func main() {
	// "urn:schemas-upnp-org:device:InternetGatewayDevice:1"
	devices, err := goupnp.DiscoverDevices(ssdp.UPNPRootDevice /*ssdp.SSDPAll*/)

	if err != nil {
		log.WithError(err).Fatal("failed to discover devices")
	}

	for _, device := range devices {
		if device.Err != nil {
			log.WithError(device.Err).WithField("url", device.Location.String()).Warning("failed to discover device. skipping.")
			continue
		}

		fmt.Printf("%s\n", device.Location.Host)

		fmt.Printf("    Devices\n")
		device.Root.Device.VisitDevices(func(dev *goupnp.Device) {
			fmt.Printf("        %s - %s - %s\n", dev.FriendlyName, dev.Manufacturer, dev.DeviceType)
		})

		fmt.Printf("    Services:\n")
		device.Root.Device.VisitServices(func(svc *goupnp.Service) {
			fmt.Printf("        %s\n", svc.String())

		})
	}
}
