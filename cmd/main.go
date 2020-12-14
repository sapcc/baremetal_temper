package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sapcc/ironic_temper/pkg/config"
	"github.com/sapcc/ironic_temper/pkg/provision"
	log "github.com/sirupsen/logrus"
)

var errors chan error
var opts config.Options

func init() {
	// default log level
	opts.LogLevelValue = config.LogLevelValue{LogLevel: log.InfoLevel}

	flag.StringVar(&opts.ConfigFilePath, "CONFIG_FILE", "./etc/config.yaml", "Path to the config file")
	flag.DurationVar(&opts.CheckInterval, "CHECK_INTERVAL", 60*time.Minute, "interval for the check")
	flag.Var(&opts.LogLevelValue, "LOG_LEVEL", "log level")
	flag.Parse()

	log.SetLevel(opts.LogLevelValue.LogLevel)
}

func main() {
	c := make(chan os.Signal, 1)
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
	cfg, err := config.GetConfig(opts)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	r := provision.NewScheduler(ctx, cfg)
	r.Start(opts.CheckInterval)
}
