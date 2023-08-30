// The config package is responsible for setting up the configuration file.
package config

import (
	"log"
	"sync"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	HTTP_port   int           `env:"HTTP_PORT"    envDefault:"3000"`
	IsProd      bool          `env:"IS_PROD"      envDefault:"false"`
	DB_URL      string        `env:"DB_URL"       envDefault:"postgres://postgres:qwerty@service-db:5432/segmentation?sslmode=disable"`
	Timeout     time.Duration `env:"TIMEOUT"      envDefault:"10s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
}

var (
	config Config = Config{}
	once   sync.Once
)

// Get returns a config structure initialized with the values ​​of the environment variables.
func Get() *Config {
	once.Do(func() {
		if err := env.Parse(&config); err != nil {
			log.Fatalf("getting config failed: %s", err.Error())
		}
	})
	return &config
}
