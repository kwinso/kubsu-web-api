package middleware

import (
	"net/http"
	"strings"
)

type mimeTypeEnforceMiddleware struct {
	mimeTypes []string
	handler   http.Handler
}

func NewMimeTypeEnforce(handler http.Handler, mimeTypes ...string) http.Handler {
	return &mimeTypeEnforceMiddleware{
		mimeTypes: mimeTypes,
		handler:   handler,
	}
}

func (m mimeTypeEnforceMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "GET" {
		m.handler.ServeHTTP(w, r)
		return
	}

	for _, mimeType := range m.mimeTypes {
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, mimeType) {
			m.handler.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
}
