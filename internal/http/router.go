package http

import (
	"net/http"
	"strings"
)

func Register(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", h.Shorten)

	mux.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/r/")
		h.GetOriginalURL(w, r, id)
	})

	return mux
}
