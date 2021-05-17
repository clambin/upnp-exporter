package main

import (
	"fmt"
	"github.com/clambin/upnp-exporter/internal/upnpstats"
	"github.com/clambin/upnp-exporter/internal/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	Port     int
	Debug    bool
	Interval time.Duration
	URL      *url.URL
	Discover bool
)

func parseOptions() {
	var urlAsStr string
	var err error

	a := kingpin.New(filepath.Base(os.Args[0]), "upnp-exporter")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&Debug)
	a.Flag("port", "Prometheus listener port").Short('p').Default("8080").IntVar(&Port)
	a.Flag("interval", "Measurement interval").Short('i').Default("30s").DurationVar(&Interval)
	a.Flag("url", "Router Service URL").Short('u').StringVar(&urlAsStr)
	a.Flag("discover", "Discover router URLs and exit").BoolVar(&Discover)

	_, err = a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if urlAsStr != "" {
		URL, err = url.Parse(urlAsStr)
		if err != nil {
			log.WithError(err).WithField("url", urlAsStr).Fatal("unable to parse router URL")
		}
	}

	if Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func discoverURLS() {
	routers, err := upnpstats.DiscoverRouters()

	if err != nil {
		log.WithError(err).Fatal("unable to discover router URLs")
	}

	for _, router := range routers {
		fmt.Printf("router found: %s - %s\n", router.URLBaseStr, router.Device.FriendlyName)
	}
}

func main() {
	parseOptions()

	log.WithField("version", version.BuildVersion).Info("upnp-exporter started")

	if Discover {
		discoverURLS()
		os.Exit(0)
	}

	go upnpstats.Run(URL, Interval)

	// Run initialized & runs the metrics
	listenAddress := fmt.Sprintf(":%d", Port)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddress, nil)
	log.WithError(err).Fatal("Failed to start Prometheus http handler")
}
