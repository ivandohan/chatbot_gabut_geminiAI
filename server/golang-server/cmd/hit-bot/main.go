package main

import (
	"golang-server/internal/boot"
	"log"
)

func main() {
	log.Println("http://localhost:8080/dohan")

	if err := boot.HTTP(); err != nil {
		log.Println("[HTTP] failed to boot API server due to error,", err.Error())
	}
}
