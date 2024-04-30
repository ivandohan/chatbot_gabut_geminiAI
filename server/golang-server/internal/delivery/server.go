package delivery

import (
	"context"
	"errors"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ChatbotHandler interface {
	ChatbotHandler(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	ChatBot ChatbotHandler
}

func (s *Server) Serve(port string) error {
	handler := cors.AllowAll().Handler(s.Handler())
	return ServeGracefully(port, handler)
}

func ServeGracefully(port string, h http.Handler) error {
	server := &http.Server{
		ReadTimeout:  1000 * time.Second,
		WriteTimeout: 1000 * time.Second,
		Handler:      h,
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	idleConnectionsClosed := make(chan os.Signal, 1)
	go func() {
		signals := make(chan os.Signal, 1)

		signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		<-signals

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("API server shutdown, error : %v", err)
		}

		close(idleConnectionsClosed)
	}()

	log.Println("API server running on port", port)
	if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-idleConnectionsClosed
	log.Println("API server shutdown gracefully.")

	return nil
}
