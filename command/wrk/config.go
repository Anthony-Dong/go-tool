package wrk

import (
	"net/http"
	"runtime"
	"time"
)

var (
	defaultConfig = Config{
		Threads:     runtime.NumCPU(),
		Connections: runtime.NumCPU(),
		Timeout:     time.Second * 5,
		Method:      http.MethodGet,
	}
)

type Option func(config *Config)

func initOp(config *Config, option ...Option) {
	for _, elem := range option {
		elem(config)
	}
	initConfig(config)
}

func initConfig(config *Config) {
	if config.Connections <= 0 {
		config.Connections = defaultConfig.Connections
	}
	if config.Threads <= 0 {
		config.Threads = defaultConfig.Threads
	}
	if config.Method == "" {
		config.Method = defaultConfig.Method
	}
}
