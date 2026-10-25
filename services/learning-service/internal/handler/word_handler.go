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

// Register はルートを登録する。どれもログインが必要で、登録・更新・削除は管理者のみ（service で確かめる）
func (h *WordHandler) Register(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler) {
	handle := func(pattern string, f http.HandlerFunc) {
		mux.Handle(pattern, requireAuth(f))
	}
	handle("GET /words", h.List)
	handle("POST /words", h.Create)
	handle("GET /words/{id}", h.Get)
	handle("PUT /words/{id}", h.Update)
	handle("DELETE /words/{id}", h.Delete)
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
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	var in model.WordInput
	if !decodeJSON(w, r, &in) {
		return
	}
	word, err := h.service.CreateWord(r.Context(), actor, in)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/words/%d", word.ID))
	writeJSON(w, http.StatusCreated, word)
}

// Get は GET /words/{id}
func (h *WordHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "id")
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
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	var in model.WordInput
	if !decodeJSON(w, r, &in) {
		return
	}
	word, err := h.service.UpdateWord(r.Context(), actor, id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, word)
}

// Delete は DELETE /words/{id}
func (h *WordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	if err := h.service.DeleteWord(r.Context(), actor, id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
