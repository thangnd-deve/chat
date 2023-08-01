package main

import (
	"chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	routes := routes()
	log.Println("Service running 8000")
	go handlers.ListenToWsChannel()
	_ = http.ListenAndServe(":8000", routes)
}
