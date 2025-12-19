package handler

import (
	"net/http"

	"github.com/kwun/onelink-api/lib"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	lib.Handler(w, r)
}
