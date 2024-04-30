package chatbot

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"golang-server/pkg/gemini_helper"
	"log"
)

type Data struct {
	genaiClient *genai.Client
}

func New(aiClient *genai.Client) Data {
	return Data{
		genaiClient: aiClient,
	}
}

func (d Data) BotGreetings(ctx context.Context, clientText string) (resultString string, err error) {

	log.Println("[Layer - Data] BotGreetings data is hit.")

	aiModel := d.genaiClient.GenerativeModel("gemini-pro")

	aiCS := aiModel.StartChat()

	response, err := aiCS.SendMessage(ctx, genai.Text(clientText))

	resultString, err = gemini_helper.ParsingTextResponseForChat(response)

	log.Println("RESULT: ", resultString)

	return resultString, err
}
