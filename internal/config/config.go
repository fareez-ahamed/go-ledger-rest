package config

import (
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string
}

func Load() (*Config, error) {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Host: host,
		Port: port,
	}, nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
