package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type ChatRequest struct {
	Message string `json:"message"`
	Profile string `json:"profile"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// respondWithError sends a JSON error response
// respondWithError发送JSON错误响应
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for frontend requests
	// 为前端请求启用CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	// 处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		respondWithError(w, "OPENROUTER_API_KEY environment variable is not configured", http.StatusInternalServerError)
		return
	}

	// Prepare OpenRouter request
	// 准备OpenRouter请求
	openRouterReqBody := map[string]interface{}{
		"model": "nvidia/nemotron-3-nano-30b-a3b:free",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an AI assistant for Kwun Tools, a collection of web-based utilities. You help users with various tasks including image analysis, JSON formatting, QR code generation, Base64 encoding/decoding, and other web tools. Be helpful, concise, and friendly. Provide clear explanations and step-by-step guidance when needed.",
			},
			{
				"role":    "user",
				"content": req.Message,
			},
		},
	}

	jsonData, err := json.Marshal(openRouterReqBody)
	if err != nil {
		respondWithError(w, "Failed to prepare request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		respondWithError(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://kwun-tools.vercel.app")
	httpReq.Header.Set("X-Title", "Kwun Tools")

	resp, err := client.Do(httpReq)
	if err != nil {
		respondWithError(w, "Failed to connect to OpenRouter: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if OpenRouter returned an error
	// 检查OpenRouter是否返回错误
	if resp.StatusCode != http.StatusOK {
		respondWithError(w, "OpenRouter API error ("+resp.Status+"): "+string(body), http.StatusBadGateway)
		return
	}

	var openRouterResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		respondWithError(w, "Failed to parse OpenRouter response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(openRouterResp.Choices) == 0 {
		respondWithError(w, "OpenRouter returned no response choices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{
		Reply: openRouterResp.Choices[0].Message.Content,
	})
}
