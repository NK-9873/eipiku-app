package handler

import (
	"fmt"
	"net/http"

	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/service"
)

type CourseHandler struct {
	service *service.CourseService
}

func NewCourseHandler(s *service.CourseService) *CourseHandler {
	return &CourseHandler{service: s}
}

// Register は /courses 以下のルートを登録する
func (h *CourseHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /courses", h.List)
	mux.HandleFunc("POST /courses", h.Create)
	mux.HandleFunc("GET /courses/{id}", h.Get)
	mux.HandleFunc("PUT /courses/{id}", h.Update)
	mux.HandleFunc("DELETE /courses/{id}", h.Delete)
	mux.HandleFunc("GET /courses/{id}/words", h.ListWords)
	mux.HandleFunc("PUT /courses/{id}/words", h.ReplaceWords)
}

// List は GET /courses
func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.ListCourses(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, courses)
}

// Create は POST /courses
func (h *CourseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in model.CourseInput
	if !decodeJSON(w, r, &in) {
		return
	}
	course, err := h.service.CreateCourse(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/courses/%d", course.ID))
	writeJSON(w, http.StatusCreated, course)
}

// Get は GET /courses/{id}
func (h *CourseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	course, err := h.service.GetCourse(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, course)
}

// Update は PUT /courses/{id}
func (h *CourseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in model.CourseInput
	if !decodeJSON(w, r, &in) {
		return
	}
	course, err := h.service.UpdateCourse(r.Context(), id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, course)
}

// Delete は DELETE /courses/{id}
func (h *CourseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.service.DeleteCourse(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListWords は GET /courses/{id}/words
func (h *CourseHandler) ListWords(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	words, err := h.service.ListCourseWords(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, words)
}

// ReplaceWords は PUT /courses/{id}/words。body は {"word_ids": [3, 1, 2]}（この順が出題順）
func (h *CourseHandler) ReplaceWords(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in model.CourseWordsInput
	if !decodeJSON(w, r, &in) {
		return
	}
	words, err := h.service.ReplaceCourseWords(r.Context(), id, in.WordIDs)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, words)
}
