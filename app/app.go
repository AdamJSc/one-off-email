package app

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Config stores our app's config
type Config struct {
	MailgunAPIKey string `envconfig:"MAILGUN_API_KEY"`
}

// MustParseConfig returns an inflated Config object from the provided file path
// or fails on error
func MustParseConfig(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal(err)
	}

	var cfg Config
	envconfig.MustProcess("", &cfg)

	return &cfg
}

// Dependencies defines the interface for our app's dependencies
type Dependencies interface {
	Config() *Config
}
