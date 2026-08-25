package source

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type handler struct {
	store Store
}

func (h *handler) AddSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req rssSourceCreateDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusBadRequest)
		return
	}
	_, err := h.store.AddSource(ctx, rssSource{
		Id:   0,
		Name: req.Name,
		Url:  req.Url,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) GetSources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sources, err := h.store.GetSources(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}
	var response []rssSourceResponseDto
	for _, source := range sources {
		response = append(response, rssSourceResponseDto{source.Id, source.Name, source.Url})
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}
}
