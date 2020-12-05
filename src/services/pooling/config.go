package pooling

import (
	"os"
)

// Config holds configuration data.
type Config struct {
	Stage     string `env:"STAGE"`
	Region    string `env:"REGION"`
	AccountID string `env:"ACCOUNT_ID"`
}

// GetConfig gets the config.
func GetConfig() *Config {
	config := Config{
		Stage:     os.Getenv("STAGE"),
		Region:    os.Getenv("REGION"),
		AccountID: os.Getenv("ACCOUNT_ID"),
	}
	return &config
}
