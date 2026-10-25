package handler

import (
	"net/http"
	"strconv"

	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/service"
)

// ProgressHandler は /me 以下（ログイン中のユーザー自身の学習記録）と、管理者用の学習記録の参照
type ProgressHandler struct {
	service *service.ProgressService
}

func NewProgressHandler(s *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{service: s}
}

// Register は /me 以下と管理者用のルートを登録する。requireAuth は JWT を検証してユーザー情報を context に入れる middleware
func (h *ProgressHandler) Register(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler) {
	handle := func(pattern string, f http.HandlerFunc) {
		mux.Handle(pattern, requireAuth(f))
	}
	handle("GET /me/word-progresses", h.ListWordProgresses)
	handle("PUT /me/word-progresses/{wordId}", h.SaveWordProgress)
	handle("GET /me/review-words", h.ListReviewWords)
	handle("GET /me/course-progresses", h.ListCourseProgresses)
	handle("PUT /me/course-progresses/{courseId}", h.SaveCourseProgress)
	handle("GET /me/summary", h.Summary)
	// 管理者かどうかは service で確かめる
	handle("GET /admin/users/{id}/progress", h.UserProgress)
}

// ListWordProgresses は GET /me/word-progresses?course_id=1&mastery=weak,unsure
func (h *ProgressHandler) ListWordProgresses(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	q, ok := progressQuery(w, r)
	if !ok {
		return
	}
	progresses, err := h.service.ListWordProgresses(r.Context(), actor.UserID, q)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progresses)
}

// SaveWordProgress は PUT /me/word-progresses/{wordId}。body は {"mastery": "weak"}
func (h *ProgressHandler) SaveWordProgress(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	wordID, ok := pathID(w, r, "wordId")
	if !ok {
		return
	}
	var in model.WordProgressInput
	if !decodeJSON(w, r, &in) {
		return
	}
	progress, err := h.service.SaveWordProgress(r.Context(), actor.UserID, wordID, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progress)
}

// ListReviewWords は GET /me/review-words?course_id=1&mastery=weak
func (h *ProgressHandler) ListReviewWords(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	q, ok := progressQuery(w, r)
	if !ok {
		return
	}
	words, err := h.service.ListReviewWords(r.Context(), actor.UserID, q)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, words)
}

// ListCourseProgresses は GET /me/course-progresses
func (h *ProgressHandler) ListCourseProgresses(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	progresses, err := h.service.ListCourseProgresses(r.Context(), actor.UserID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progresses)
}

// SaveCourseProgress は PUT /me/course-progresses/{courseId}。body は {"current_position": 2}
func (h *ProgressHandler) SaveCourseProgress(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	courseID, ok := pathID(w, r, "courseId")
	if !ok {
		return
	}
	var in model.CourseProgressInput
	if !decodeJSON(w, r, &in) {
		return
	}
	progress, err := h.service.SaveCourseProgress(r.Context(), actor.UserID, courseID, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progress)
}

// Summary は GET /me/summary
func (h *ProgressHandler) Summary(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	summary, err := h.service.Summary(r.Context(), actor.UserID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// UserProgress は GET /admin/users/{id}/progress（管理者のみ）
func (h *ProgressHandler) UserProgress(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	userID, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	progress, err := h.service.UserProgress(r.Context(), actor, userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progress)
}

// progressQuery はクエリの course_id と mastery を読む。course_id が数字でなければ 400
func progressQuery(w http.ResponseWriter, r *http.Request) (service.ProgressQuery, bool) {
	q := service.ProgressQuery{Mastery: r.URL.Query().Get("mastery")}
	if raw := r.URL.Query().Get("course_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || id == 0 {
			writeBadRequest(w, "course_id が不正です")
			return service.ProgressQuery{}, false
		}
		q.CourseID = &id
	}
	return q, true
}
