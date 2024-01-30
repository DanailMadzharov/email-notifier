package message

import (
	"context"
	"github.com/jdvr/go-again"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type ProduceRetryEventOperation struct {
	writer  *kafka.Writer
	message *[]byte
}

func Produce(ctx context.Context, message *[]byte, producer *kafka.Writer,
	config *again.BackoffConfiguration) {
	retryService := again.WithExponentialBackoff[bool](*config)
	retryOperation := ProduceRetryEventOperation{
		writer:  producer,
		message: message,
	}
	_, err := retryService.Retry(ctx, &retryOperation)
	if err != nil {
		log.Error().Msgf("Failed to send object in dead letter queue: %s", err.Error())
	}

	log.Info().Msg("Successfully sent object to dead letter topic")
}

func (operation *ProduceRetryEventOperation) Run(ctx context.Context) (bool, error) {
	err := operation.writer.WriteMessages(ctx, kafka.Message{
		Value: *operation.message,
	})
	if err != nil {
		return true, err
	}

	return true, nil
}
