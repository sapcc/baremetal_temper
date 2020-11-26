package config

import "time"

// Options passed via cmd line
type Options struct {
	LogLevel       string
	Version        string
	ConfigFilePath string
	CheckInterval  time.Duration
	NameSpace      string
	Region         string
}
