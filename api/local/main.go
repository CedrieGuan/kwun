package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kwun/onelink-api/lib"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found in .env, checking standard environment")
	}

	http.HandleFunc("/api/chat", lib.Handler)

	log.Println("Local server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
