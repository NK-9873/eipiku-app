package handler

import (
	"fmt"
	"net/http"

	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/service"
)

type WordHandler struct {
	service *service.WordService
}

func NewWordHandler(s *service.WordService) *WordHandler {
	return &WordHandler{service: s}
}

// Register は /words 以下のルートを登録する
func (h *WordHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /words", h.List)
	mux.HandleFunc("POST /words", h.Create)
	mux.HandleFunc("GET /words/{id}", h.Get)
	mux.HandleFunc("PUT /words/{id}", h.Update)
	mux.HandleFunc("DELETE /words/{id}", h.Delete)
}

// List は GET /words
func (h *WordHandler) List(w http.ResponseWriter, r *http.Request) {
	words, err := h.service.ListWords(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, words)
}

// Create は POST /words
func (h *WordHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in model.WordInput
	if !decodeJSON(w, r, &in) {
		return
	}
	word, err := h.service.CreateWord(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/words/%d", word.ID))
	writeJSON(w, http.StatusCreated, word)
}

// Get は GET /words/{id}
func (h *WordHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	word, err := h.service.GetWord(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, word)
}

// Update は PUT /words/{id}
func (h *WordHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in model.WordInput
	if !decodeJSON(w, r, &in) {
		return
	}
	word, err := h.service.UpdateWord(r.Context(), id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, word)
}

// Delete は DELETE /words/{id}
func (h *WordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.service.DeleteWord(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
