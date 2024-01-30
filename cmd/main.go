package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"sumup-email-notifier/config"
	"sumup-email-notifier/message"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal().Msgf("%s", err.Error())
	}
	log.Info().Msgf("Application Version: %s", viper.GetString("version"))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go message.ConsumeRecoveryEvents(ctx)
	message.Consume(ctx)
}
