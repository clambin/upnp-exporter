package main

import (
	"fmt"
	"github.com/huin/goupnp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"upnp-exporter/internal/upnpstats"
	"upnp-exporter/internal/version"
)

var (
	Port     int
	Debug    bool
	Interval time.Duration
)

func parseOptions() {
	a := kingpin.New(filepath.Base(os.Args[0]), "upnp-exporter")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&Debug)
	a.Flag("port", "Prometheus listener port").Short('p').Default("8080").IntVar(&Port)
	a.Flag("interval", "Measurement interval").Short('i').Default("30s").DurationVar(&Interval)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func scrape(routers []*goupnp.RootDevice) {
	for {
		for _, router := range routers {
			upnpstats.ReportNetworkStats(router)
		}
		time.Sleep(Interval)
	}
}

func main() {
	parseOptions()

	log.WithField("version", version.BuildVersion).Info("upnp-exporter started")

	routers, err := upnpstats.DiscoverRouters()

	if err != nil {
		log.WithField("err", err).Fatal("unable to discover routers. exiting")
	}

	go scrape(routers)

	// Run initialized & runs the metrics
	listenAddress := fmt.Sprintf(":%d", Port)
	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(listenAddress, nil)
	log.WithError(err).Fatal("Failed to start Prometheus http handler")
}
