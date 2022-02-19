package app

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS,notEmpty" envDefault:"localhost:8080"`
	BaseURL string `env:"BASE_URL,notEmpty" envDefault:"http://localhost:8080"`
}

func (c *ServerConfig) LoadEnvs() {
	if err := env.Parse(c); err != nil {
		log.Printf("%+v\n", err)
	}
}
