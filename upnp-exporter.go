package main

import (
	"fmt"
	"github.com/clambin/gotools/metrics"
	"github.com/clambin/upnp-exporter/collector"
	"github.com/clambin/upnp-exporter/upnpstats"
	"github.com/clambin/upnp-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
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

func main() {
	parseOptions()

	scanner, err := upnpstats.New(URL)
	if err != nil {
		log.WithError(err).Fatal("unable to create upnp scanner")
	}

	if Discover {
		for index, router := range scanner.Routers() {
			fmt.Printf("%d: %s\n", index+1, router)
		}
		os.Exit(0)
	}

	log.WithField("version", version.BuildVersion).Info("upnp-exporter started")
	c := collector.New(scanner)
	prometheus.MustRegister(c)

	// Run initialized & runs the metrics
	server := metrics.NewServer(Port)
	go func() {
		err2 := server.Run()
		if err2 != nil && err2 != http.ErrServerClosed {
			log.WithError(err2).Fatal("Failed to start Prometheus http handler")
		}
	}()

	<-shutdown.Chan()
	err = server.Shutdown(30 * time.Second)
	if err != nil {
		log.WithError(err).Fatal("failed to do graceful shutdown for given time")
	}
	log.Info("upnp-exporter stopped")
}
