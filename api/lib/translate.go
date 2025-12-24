package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// TranslateRequest represents the translation request from the frontend
// TranslateRequest表示来自前端的翻译请求
type TranslateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"` // empty string for auto-detect / 空字符串表示自动检测
	TargetLang string `json:"target_lang"`
}

// TranslateResponse represents the translation response to the frontend
// TranslateResponse表示返回给前端的翻译响应
type TranslateResponse struct {
	TranslatedText   string `json:"translated_text"`
	DetectedLanguage string `json:"detected_language,omitempty"`
}

// DeepLTranslateResponse represents the response from DeepL API
// DeepLTranslateResponse表示DeepL API的响应
type DeepLTranslateResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

// TranslateHandler handles translation requests using DeepL API
// TranslateHandler使用DeepL API处理翻译请求
func TranslateHandler(w http.ResponseWriter, r *http.Request) {
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

	var req TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request parameters
	// 验证请求参数
	if req.Text == "" {
		respondWithError(w, "Text is required", http.StatusBadRequest)
		return
	}

	if req.TargetLang == "" {
		respondWithError(w, "Target language is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("DEEPL_API_KEY")
	if apiKey == "" {
		respondWithError(w, "DEEPL_API_KEY environment variable is not configured", http.StatusInternalServerError)
		return
	}

	// Prepare DeepL API request
	// 准备DeepL API请求
	deeplURL := "https://api-free.deepl.com/v2/translate"

	// Build form data
	// 构建表单数据
	formData := url.Values{}
	formData.Set("text", req.Text)
	formData.Set("target_lang", strings.ToUpper(req.TargetLang))
	if req.SourceLang != "" && req.SourceLang != "auto" {
		formData.Set("source_lang", strings.ToUpper(req.SourceLang))
	}

	// Create HTTP request
	// 创建HTTP请求
	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", deeplURL, strings.NewReader(formData.Encode()))
	if err != nil {
		respondWithError(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
		respondWithError(w, fmt.Sprintf("DeepL API error (%s): %s", resp.Status, string(body)), http.StatusBadGateway)
		return
	}

	// Parse DeepL response
	// 解析DeepL响应
	var deeplResp DeepLTranslateResponse
	if err := json.Unmarshal(body, &deeplResp); err != nil {
		respondWithError(w, "Failed to parse DeepL response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(deeplResp.Translations) == 0 {
		respondWithError(w, "DeepL returned no translations", http.StatusInternalServerError)
		return
	}

	// Build response
	// 构建响应
	translationResult := TranslateResponse{
		TranslatedText:   deeplResp.Translations[0].Text,
		DetectedLanguage: deeplResp.Translations[0].DetectedSourceLanguage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(translationResult)
}
