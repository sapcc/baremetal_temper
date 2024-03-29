/**
 * Copyright 2021 SAP SE
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sapcc/baremetal_temper/cmd"
	"github.com/sapcc/baremetal_temper/pkg/config"
	"github.com/sapcc/baremetal_temper/pkg/server"
	"github.com/sapcc/baremetal_temper/pkg/temper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var opts config.Options
var wait time.Duration

func main() {
	var cfg config.Config
	viper.SetConfigFile(opts.ConfigFilePath)
	cmd.InitConfig()
	cmd.UnmarshalConfig(&cfg)
	ctxLogger := log.WithFields(log.Fields{
		"temper": "server",
	})
	t := temper.New(opts.Workers)
	s := server.New(cfg, ctxLogger, t)
	s.RegisterAPIRoutes()
	srv := &http.Server{
		Addr: "0.0.0.0:80",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		// https://operations.global.cloud.sap/docs/support/playbook/kubernetes/idle_http_keep_alive_timeout.html
		ReadTimeout: time.Second * 61,
		IdleTimeout: time.Second * 61,
		Handler:     s.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}

func init() {
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.IntVar(&opts.Workers, "number-workers", 10, "set the number of temper workers to handle tempering of nodes simultaneously")
	flag.StringVar(&opts.ConfigFilePath, "config-path", "etc/config/temper.yaml", "set the path to the config file")
	flag.Parse()

	// default log level
	opts.LogLevelValue = config.LogLevelValue{LogLevel: log.DebugLevel}
	log.SetLevel(opts.LogLevelValue.LogLevel)

	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		hook, err := logrus_sentry.NewSentryHook(dsn, []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})
		if err == nil {
			log.Info("adding sentry hook")
			log.AddHook(hook)
		}
	}
}
