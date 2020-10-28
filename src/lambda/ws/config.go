package lambdaws

import (
	"fmt"
	"os"

	"github.com/tkanos/gonfig"
)

// Config holds configuration data.
type Config struct {
	MemcachedHost     string
	MemcachedUsername string
	MemcachedPassword string
	Stage             string `env:"STAGE"`
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
