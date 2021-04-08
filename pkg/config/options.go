package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Options passed via cmd line
type Options struct {
	LogLevelValue  LogLevelValue
	Baremetal      bool
	Diagnostics    bool
	Version        string
	ConfigFilePath string
	RulesFilePath  string
	CheckInterval  time.Duration
}

type LogLevelValue struct {
	LogLevel log.Level
}

func (l LogLevelValue) String() string {
	return l.LogLevel.String()
}

func (l LogLevelValue) Set(s string) error {
	if logLevel, err := log.ParseLevel(s); err != nil {
		return err
	} else {
		l.LogLevel = logLevel
	}
	return nil
}
