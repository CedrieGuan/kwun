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
	// Load API keys from environment variables
	// 从环境变量加载API密钥
	openRouterKey := os.Getenv("OPENROUTER_API_KEY")
	if openRouterKey == "" {
		log.Println("⚠️  Warning: OPENROUTER_API_KEY is not set. Chat API will not work.")
	} else {
		log.Println("✅ OPENROUTER_API_KEY is configured")
	}

	deepLKey := os.Getenv("DEEPL_API_KEY")
	if deepLKey == "" {
		log.Println("⚠️  Warning: DEEPL_API_KEY is not set. Translation API will not work.")
	} else {
		log.Println("✅ DEEPL_API_KEY is configured")
	}

	// Register API endpoints
	// 注册API端点
	http.HandleFunc("/api/chat", lib.Handler)
	http.HandleFunc("/api/translate", lib.TranslateHandler)
	http.HandleFunc("/api/usage", lib.UsageHandler)

	// Start the local development server
	// 启动本地开发服务器
	port := ":8080"
	log.Printf("🚀 Local Go server starting on http://localhost%s", port)
	log.Printf("📡 Chat API available at http://localhost%s/api/chat", port)
	log.Printf("🌐 Translate API available at http://localhost%s/api/translate", port)
	log.Printf("📊 Usage API available at http://localhost%s/api/usage", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
