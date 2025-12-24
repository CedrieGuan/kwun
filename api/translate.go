package handler

import (
	"net/http"

	"github.com/kwun/onelink-api/lib"
)

// Handler processes translation requests for Vercel serverless deployment
// Handler处理Vercel无服务器部署的翻译请求
func Handler(w http.ResponseWriter, r *http.Request) {
	lib.TranslateHandler(w, r)
}
