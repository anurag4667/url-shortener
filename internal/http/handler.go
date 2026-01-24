package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anurag4667/url-shortener/internal/kafka/producer"
	"github.com/anurag4667/url-shortener/internal/redis"
	"github.com/anurag4667/url-shortener/internal/service"
)

type Handler struct {
	service  *service.URLService
	producer *producer.ClickProducer
}

func New(service *service.URLService, producer *producer.ClickProducer) *Handler {
	return &Handler{service: service, producer: producer}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	log.Println("Shorten called") // ADD THIS

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log.Println("Decoded URL:", req.URL) // ADD THIS

	id, err := h.service.Shorten(req.URL)
	if err != nil {
		http.Error(w, "failed to shorten", http.StatusInternalServerError)
		return
	}

	log.Println("Generated ID:", id) // ADD THIS

	json.NewEncoder(w).Encode(map[string]string{
		"short_url": "http://localhost:4000/r/" + id,
	})
}

func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("GetOriginalURL called with ID:", id)

	cachedURL, err := redis.GetURL(id)

	if err == nil {
		fmt.Println("Cache hit:", cachedURL)

		go h.producer.Publish(context.Background(), id, id)

		json.NewEncoder(w).Encode(map[string]string{
			"original_url": cachedURL,
			"source":       "cache",
		})
		return
	}

	url, ok, err := h.service.Resolve(id)
	if err != nil {
		fmt.Println("DB resolve error:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid url",
		})
		return
	}

	_ = redis.SetURL(id, url)

	go h.producer.Publish(context.Background(), id, id)

	json.NewEncoder(w).Encode(map[string]string{
		"original_url": url,
		"source":       "database",
	})
}
