package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/redfish"
	log "github.com/sirupsen/logrus"
)

var errors chan error
var opts config.Options

func init() {
	flag.StringVar(&opts.ConfigFilePath, "CONFIG_FILE", "./etc/config.yaml", "Path to the config file")
	flag.DurationVar(&opts.CheckInterval, "CHECK_INTERVAL", 60*time.Second, "interval for the check")
	flag.Parse()
}

func main() {
	c := make(chan os.Signal, 1)
	errors := make(chan error, 0)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range c {
			log.Error(sig)
			cancel()
			os.Exit(0)
		}
	}()
	go func() {
		select {
		case err := <-errors:
			log.Error(err.Error())
			os.Exit(1)
		}
	}()
	cfg, err := config.GetConfig(opts)
	if err != nil {
		os.Exit(0)
	}
	r := redfish.New(cfg)
	r.Start(ctx, errors)
}
