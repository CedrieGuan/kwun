package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kwun/onelink-api/lib"
)

// main function for local development server
// 本地开发服务器的主函数
func main() {
	// Load API key from environment variable
	// 从环境变量加载API密钥
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: OPENROUTER_API_KEY environment variable is not set. Please set it before running the server.")
	}

	// Register the chat endpoint
	// 注册聊天端点
	http.HandleFunc("/api/chat", lib.Handler)

	// Start the local development server
	// 启动本地开发服务器
	port := ":8080"
	log.Printf("🚀 Local Go server starting on http://localhost%s", port)
	log.Printf("📡 Chat API available at http://localhost%s/api/chat", port)
	log.Println("✅ OPENROUTER_API_KEY is configured")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
