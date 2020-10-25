package lambdaws

import (
	"fmt"
	"os"

	"github.com/YouJinTou/vocabrace/pool"
	"github.com/tkanos/gonfig"
)

// GetPoolConfig gets the pool config.
func GetPoolConfig() *pool.Config {
	config := pool.Config{}
	stage := os.Getenv("STAGE")
	file := fmt.Sprintf("config.%s.json", stage)
	err := gonfig.GetConf(file, &config)

	if err != nil {
		panic(err)
	}

	return &config
}
