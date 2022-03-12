package app

import (
	"log"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

type ServerConfig struct {
	Host            string `env:"HOST" envDefault:"localhost"`
	Address         string `env:"SERVER_ADDRESS,notEmpty" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL,notEmpty" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
}

func (c *ServerConfig) LoadEnvs() {
	if err := env.Parse(c); err != nil {
		log.Printf("%+v\n", err)
	}
}

func (c *ServerConfig) ParseCommandLine() {
	if flag.Lookup("a") == nil {
		flag.StringVarP(&c.Address, "a", "a", c.Address, "-a localhost:8080")
	}
	if flag.Lookup("b") == nil {
		flag.StringVarP(&c.BaseURL, "b", "b", c.BaseURL, "-b http://localhost:8080")
	}
	if flag.Lookup("f") == nil {
		flag.StringVarP(&c.FileStoragePath, "f", "f", c.FileStoragePath, "-f /tmp/filename.tmp")
	}
	if flag.Lookup("database-dsn") == nil {
		flag.StringVar(&c.DatabaseDsn, "database-dsn", c.DatabaseDsn, "--database-dsn db_driver://user:pass@domain:port/db_name")
	}
	if flag.Lookup("d") == nil {
		flag.StringVarP(&c.DatabaseDsn, "d", "d", c.DatabaseDsn, "-d db_driver://user:pass@domain:port/db_name")
	}

	flag.Parse()
}
