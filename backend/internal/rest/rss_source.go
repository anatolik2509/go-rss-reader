package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"rss-reader-backend/internal/service/source"
)

type rssSourceCreateDto struct {
	Name string
	Url  string
}

type rssSourceResponseDto struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type RssSourceHandler struct {
	sourceService *source.Service
}

func (h *RssSourceHandler) AddSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req rssSourceCreateDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, "handling AddSource request failed", "error", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	sourceUrl, err := url.Parse(req.Url)
	if err != nil {
		slog.ErrorContext(ctx, "handling AddSource request failed", "error", err)
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	_, err = h.sourceService.AddSource(ctx, source.RssSource{
		Id:   0,
		Name: req.Name,
		Url:  sourceUrl,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "handling AddSource request failed", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *RssSourceHandler) GetSources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sources, err := h.sourceService.GetSources(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "handling GetSource request failed", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []rssSourceResponseDto
	for _, rssSource := range sources {
		response = append(response, rssSourceResponseDto{rssSource.Id, rssSource.Name, rssSource.Url.String()})
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.ErrorContext(ctx, "handling GetSource request failed", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func NewRssSourceHandler(service *source.Service) *RssSourceHandler {
	return &RssSourceHandler{service}
}
