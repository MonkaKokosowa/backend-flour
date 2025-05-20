package env

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Environment struct {
	Dialer struct {
		Server   string `env:"SMTP_SERVER,required"`
		Port     int    `env:"SMTP_PORT,required"`
		Username string `env:"SMTP_USER,required"`
		Password string `env:"SMTP_PASSWORD,required"`
		To       string `env:"SMTP_TO,required"`
	}
	WebServer struct {
		AllowedOrigins string `env:"ALLOWED_ORIGINS,required"`
		Port           int    `env:"PORT,required"`
	}
	Blog struct {
		FlatnotesURL string `env:"FLATNOTES_URL,required"`
	}
}

// GetEnv initializes environment variables from .env file (if present) and returns them
func GetEnv() (Environment, error) {
	// Attempt to load .env file, ignore error if it doesn't exist
	_ = godotenv.Load()

	var environment Environment
	if err := env.Parse(&environment); err != nil {
		return environment, err
	}
	return environment, nil
}
