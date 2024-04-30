package chatbot

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	innerLogger "golang-server/pkg/log"
	"golang-server/pkg/response"
	"log"
	"net/http"
)

type ChatbotService interface {
	BotGreetings(ctx context.Context, clientText string) (resultString string, err error)
}

type Handler struct {
	chatbotService ChatbotService
	tracer         opentracing.Tracer
	logger         innerLogger.Factory
}

func New(service ChatbotService, tracer opentracing.Tracer, logger innerLogger.Factory) *Handler {
	return &Handler{
		chatbotService: service,
		tracer:         tracer,
		logger:         logger,
	}
}

func (h *Handler) ChatbotHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           context.Context
		innerResponse *response.Response
		responseError response.CustomErrorModel
		result        interface{}
		metadata      interface{}
		err           error
	)

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := h.tracer.StartSpan("chatbotHandler", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	innerResponse = &response.Response{}
	defer innerResponse.RenderJSONResult(w, r)

	ctx = opentracing.ContextWithSpan(r.Context(), span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	if err == nil {
		parameters := r.URL.Query()

		switch r.Method {
		case http.MethodGet:
			{
				parameterLength := len(parameters)

				switch parameterLength {
				case 1:
					_, isTextExists := r.URL.Query()["text"]

					if isTextExists {
						//log.Println("Endpoint is hit by", parameterLength, " parameter/s my lord.")
						log.Println("client_request: How are you today?")
						result, err = h.chatbotService.BotGreetings(context.Background(), r.FormValue("text"))
					}
				default:
					err = errors.New("400")
				}
			}
		}
	}

	if err != nil {
		responseError = response.CustomErrorModel{
			ErrorCode:    101,
			ErrorMessage: "101 - Data Not Found",
			ErrorStatus:  true,
		}

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		innerResponse.StatusCode = responseError.ErrorCode
		innerResponse.Error = responseError

		return
	}

	innerResponse.Data = result
	innerResponse.Metadata = metadata

	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
}
