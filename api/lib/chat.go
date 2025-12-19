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

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		http.Error(w, "OPENROUTER_API_KEY is not set", http.StatusInternalServerError)
		return
	}

	// Prepare OpenRouter request
	openRouterReqBody := map[string]interface{}{
		"model": "nvidia/nemotron-3-nano-30b-a3b:free",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an AI assistant for OneLink. You help users optimize their link-in-bio profiles. Here is the current profile context: " + req.Profile,
			},
			{
				"role":    "user",
				"content": req.Message,
			},
		},
	}

	jsonData, err := json.Marshal(openRouterReqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://onelink-demo.vercel.app")
	httpReq.Header.Set("X-Title", "OneLink Demo")

	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Error parsing OpenRouter response: "+string(body), http.StatusInternalServerError)
		return
	}

	if len(openRouterResp.Choices) == 0 {
		http.Error(w, "No response from OpenRouter", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{
		Reply: openRouterResp.Choices[0].Message.Content,
	})
}
