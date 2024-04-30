package chatbot

import (
	"context"
	"github.com/pkg/errors"
	"log"
)

type Data interface {
	BotGreetings(ctx context.Context, clientText string) (resultString string, err error)
}

type Service struct {
	data Data
}

func New(data Data) Service {
	return Service{
		data: data,
	}
}

func (s Service) BotGreetings(ctx context.Context, clientText string) (resultString string, err error) {
	log.Println("[Layer - Service] BotGreetings service is hit.")

	resultString, err = s.data.BotGreetings(ctx, clientText)
	if err != nil {
		return resultString, errors.Wrap(err, "[SERVICE][BotGreetings]")
	}

	return resultString, err
}
