package main

import (
	"github.com/MonkaKokosowa/backend-flour/internal/env"
	"github.com/MonkaKokosowa/backend-flour/internal/mail"
	"github.com/MonkaKokosowa/backend-flour/internal/web"
	"github.com/rs/zerolog/log"
)

func main() {
	environment, err := env.GetEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load environment")

	}
	app := &web.App{
		Env:    &environment,
		Dialer: mail.GetDialer(environment),
	}

	web.StartWeb(*app)

}
