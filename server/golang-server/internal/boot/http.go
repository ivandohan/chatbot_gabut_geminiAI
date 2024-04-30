package boot

import (
	"context"
	"errors"
	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	chatbotD "golang-server/internal/data/chatbot"
	internalServer "golang-server/internal/delivery"
	chatbotH "golang-server/internal/delivery/chatbot"
	chatbotS "golang-server/internal/service/chatbot"
	innerLogger "golang-server/pkg/log"
	"golang-server/pkg/tracing"
	"google.golang.org/api/option"
	"io"
	"net/http"
)

func HTTP() (err error) {
	var (
		chatbotData chatbotD.Data

		chatbotService chatbotS.Service

		chatbotHandler *chatbotH.Handler

		server internalServer.Server
		logger *zap.Logger
	)

	genaiClient, err := genai.NewClient(context.Background(), option.WithAPIKey("my bad awokowak:) next update i'll hide this"))
	if err != nil {
		return err
	}

	logger, _ = zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	zapLogger := logger.With(zap.String("service", "dohan-chatbot"))
	zlogger := innerLogger.NewFactory(zapLogger)

	tracer, closer := tracing.Init("dohat-chatbot", zlogger)
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			return
		}

		err = genaiClient.Close()
		if err != nil {
			return
		}
	}(closer)

	chatbotData = chatbotD.New(genaiClient)

	chatbotService = chatbotS.New(chatbotData)

	chatbotHandler = chatbotH.New(chatbotService, tracer, zlogger)

	server = internalServer.Server{
		ChatBot: chatbotHandler,
	}

	if err := server.Serve(":8080"); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return err
}
