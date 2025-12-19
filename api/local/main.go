package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kwun/onelink-api/lib"
)

// Handler is the exported function that Vercel invokes for serverless deployment
// Handler是Vercel调用的导出函数，用于无服务器部署
func Handler(w http.ResponseWriter, r *http.Request) {
	lib.Handler(w, r)
}

// main is used for local development only
// main函数仅用于本地开发
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
