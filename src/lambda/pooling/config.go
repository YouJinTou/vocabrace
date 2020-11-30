package lambdapooling

import (
	"fmt"
	"os"

	"github.com/tkanos/gonfig"
)

// Config holds configuration data.
type Config struct {
	Stage     string `env:"STAGE"`
	Region    string `env:"REGION"`
	AccountID string `env:"ACCOUNT_ID"`
}

// GetConfig gets the config.
func GetConfig() *Config {
	config := Config{}
	stage := os.Getenv("STAGE")
	file := fmt.Sprintf("config.%s.json", stage)
	err := gonfig.GetConf(file, &config)

	if err != nil {
		panic(err)
	}

	return &config
}
