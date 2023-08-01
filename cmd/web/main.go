package main

import (
	"log"
	"net/http"
)

func main() {
	routes := routes()
	log.Println("Service running 8000")
	_ = http.ListenAndServe(":8000", routes)
}
