package config

import "time"

// Config defines the application configurations
type Config struct {
	Name          string `trim:"true"`
	Port          uint16
	LogLevel      int
	AccessLog     bool
	AwsLog        bool
	AwsRoleExpiry time.Duration
}
