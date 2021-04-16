package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/provision"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var errors chan error
var opts config.Options

func init() {
	// default log level
	opts.LogLevelValue = config.LogLevelValue{LogLevel: log.DebugLevel}

	flag.StringVar(&opts.ConfigFilePath, "CONFIG_FILE", "./etc/config.yaml", "Path to the config file")
	flag.DurationVar(&opts.CheckInterval, "CHECK_INTERVAL", 10*time.Minute, "interval for the check")
	flag.Var(&opts.LogLevelValue, "LOG_LEVEL", "log level")
	flag.BoolVar(&opts.Baremetal, "BAREMETAL_TEMPER", false, "run baremetal temper tasks")
	flag.BoolVar(&opts.Diagnostics, "DIAGNOSTICS_TEMPER", false, "run diagnostics temper tasks")
	flag.BoolVar(&opts.RedfishEvents, "SUBSCRIPE_REDFISH_EVENTS", false, "subscripe to redfish node events")
	flag.Parse()
	log.SetLevel(opts.LogLevelValue.LogLevel)

	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		hook, err := logrus_sentry.NewSentryHook(dsn, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		if err == nil {
			log.Info("adding sentry hook")
			log.AddHook(hook)
		}
	}
}

func main() {
	c := make(chan os.Signal, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range c {
			cancel()
			os.Exit(0)
		}
	}()
	cfg, err := config.GetConfig(opts)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	r, err := provision.NewScheduler(ctx, cfg, opts)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	r.Start(opts.CheckInterval)
}
