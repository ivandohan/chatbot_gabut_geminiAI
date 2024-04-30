package delivery

import (
	"errors"
	"github.com/gorilla/mux"
	response2 "golang-server/pkg/response"
	"log"
	"net/http"
)

func (s *Server) Handler() *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	router.HandleFunc("", defaultHandler).Methods("GET")
	router.HandleFunc("/", defaultHandler).Methods("GET")

	sub := router.PathPrefix("/dohan").Subrouter()

	sub.HandleFunc("", defaultHandler).Methods("GET")
	sub.HandleFunc("/", defaultHandler).Methods("GET")

	sub.HandleFunc("/chat-bot", s.ChatBot.ChatbotHandler).Methods("GET", "POST")

	return router
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	var (
		responseError response2.CustomErrorModel
		response      *response2.Response = &response2.Response{}
		innerError    error               = errors.New("404 Not Found")
	)

	defer response.RenderJSONResult(w, r)

	responseError = response2.CustomErrorModel{
		ErrorCode:    http.StatusNotFound,
		ErrorMessage: "404 Not Found",
		ErrorStatus:  true,
	}

	log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, innerError)
	response.StatusCode = http.StatusNotFound
	response.Error = responseError

	return
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("API Service for Chat Bot App."))
}
