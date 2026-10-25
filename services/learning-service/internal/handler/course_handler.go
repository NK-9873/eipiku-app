package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"eipiku/learning-service/internal/service"
)

type CourseHandler struct {
	service *service.CourseService
}

func NewCourseHandler(s *service.CourseService) *CourseHandler {
	return &CourseHandler{service: s}
}

// List は GET /courses
func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.ListCourses(r.Context())
	if err != nil {
		log.Printf("コース一覧の取得に失敗: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, courses)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("レスポンスの書き込みに失敗: %v", err)
	}
}
