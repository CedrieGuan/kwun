package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// UsageResponse represents the DeepL usage/quota response
// UsageResponse表示DeepL使用量/配额响应
type UsageResponse struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}

// DeepLUsageResponse represents the actual response from DeepL usage API
// DeepLUsageResponse表示DeepL使用量API的实际响应
type DeepLUsageResponse struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}

// UsageHandler retrieves DeepL API usage statistics
// UsageHandler检索DeepL API使用统计信息
func UsageHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for frontend requests
	// 为前端请求启用CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	// 处理预检请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("DEEPL_API_KEY")
	if apiKey == "" {
		respondWithError(w, "DEEPL_API_KEY environment variable is not configured", http.StatusInternalServerError)
		return
	}

	// Prepare DeepL usage API request
	// 准备DeepL使用量API请求
	deeplURL := "https://api-free.deepl.com/v2/usage"

	// Create HTTP request
	// 创建HTTP请求
	client := &http.Client{}
	httpReq, err := http.NewRequest("GET", deeplURL, nil)
	if err != nil {
		respondWithError(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)

	// Execute request
	// 执行请求
	resp, err := client.Do(httpReq)
	if err != nil {
		respondWithError(w, "Failed to connect to DeepL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if DeepL returned an error
	// 检查DeepL是否返回错误
	if resp.StatusCode != http.StatusOK {
		// Graceful degradation: return empty usage data instead of error
		// 优雅降级：返回空使用量数据而不是错误
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UsageResponse{
			CharacterCount: 0,
			CharacterLimit: 0,
		})
		return
	}

	// Parse DeepL response
	// 解析DeepL响应
	var deeplResp DeepLUsageResponse
	if err := json.Unmarshal(body, &deeplResp); err != nil {
		respondWithError(w, "Failed to parse DeepL response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response
	// 构建响应
	usageResult := UsageResponse{
		CharacterCount: deeplResp.CharacterCount,
		CharacterLimit: deeplResp.CharacterLimit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usageResult)
}
