package handler

import (
	"net/http"

	"github.com/kwun/onelink-api/lib"
)

// Handler processes DeepL usage/quota requests for Vercel serverless deployment
// Handler处理Vercel无服务器部署的DeepL使用量/配额请求
func Handler(w http.ResponseWriter, r *http.Request) {
	lib.UsageHandler(w, r)
}
