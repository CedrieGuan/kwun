package handler

import (
	"net/http"

	"github.com/kwun/onelink-api/lib"
)

// Handler is the exported function that Vercel invokes for serverless deployment
// Handler是Vercel调用的导出函数，用于无服务器部署
func Handler(w http.ResponseWriter, r *http.Request) {
	lib.Handler(w, r)
}
