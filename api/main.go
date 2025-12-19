package main

import (
	"log"
	"net/http"

	"github.com/kwun/onelink-api/handler"
)

func main() {
	http.HandleFunc("/api/chat", handler.ChatHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
