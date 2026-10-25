package handler

import (
	"net/http"
	"strconv"

	"eipiku/learning-service/internal/auth"
	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/service"
)

// ProgressHandler は /me 以下（ログイン中のユーザー自身の学習記録）
type ProgressHandler struct {
	service *service.ProgressService
}

func NewProgressHandler(s *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{service: s}
}

// Register は /me 以下のルートを登録する。requireUser はユーザー ID を context に入れる middleware
func (h *ProgressHandler) Register(mux *http.ServeMux, requireUser func(http.Handler) http.Handler) {
	handle := func(pattern string, f http.HandlerFunc) {
		mux.Handle(pattern, requireUser(f))
	}
	handle("GET /me/word-progresses", h.ListWordProgresses)
	handle("PUT /me/word-progresses/{wordId}", h.SaveWordProgress)
	handle("GET /me/review-words", h.ListReviewWords)
	handle("GET /me/course-progresses", h.ListCourseProgresses)
	handle("PUT /me/course-progresses/{courseId}", h.SaveCourseProgress)
	handle("GET /me/summary", h.Summary)
}

// ListWordProgresses は GET /me/word-progresses?course_id=1&mastery=weak,unsure
func (h *ProgressHandler) ListWordProgresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUser(w, r)
	if !ok {
		return
	}
	q, ok := progressQuery(w, r)
	if !ok {
		return
	}
	progresses, err := h.service.ListWordProgresses(r.Context(), userID, q)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progresses)
}

// SaveWordProgress は PUT /me/word-progresses/{wordId}。body は {"mastery": "weak"}
func (h *ProgressHandler) SaveWordProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUser(w, r)
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
	progress, err := h.service.SaveWordProgress(r.Context(), userID, wordID, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progress)
}

// ListReviewWords は GET /me/review-words?course_id=1&mastery=weak
func (h *ProgressHandler) ListReviewWords(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUser(w, r)
	if !ok {
		return
	}
	q, ok := progressQuery(w, r)
	if !ok {
		return
	}
	words, err := h.service.ListReviewWords(r.Context(), userID, q)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, words)
}

// ListCourseProgresses は GET /me/course-progresses
func (h *ProgressHandler) ListCourseProgresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUser(w, r)
	if !ok {
		return
	}
	progresses, err := h.service.ListCourseProgresses(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progresses)
}

// SaveCourseProgress は PUT /me/course-progresses/{courseId}。body は {"current_position": 2}
func (h *ProgressHandler) SaveCourseProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUser(w, r)
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
	progress, err := h.service.SaveCourseProgress(r.Context(), userID, courseID, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, progress)
}

// Summary は GET /me/summary
func (h *ProgressHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUser(w, r)
	if !ok {
		return
	}
	summary, err := h.service.Summary(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// currentUser は middleware が context に入れたユーザー ID を取り出す
func currentUser(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	userID, ok := auth.UserIDFrom(r.Context())
	if !ok {
		// middleware を通していないルートで呼ばれたとき（登録の書き忘れ）
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "ログインが必要です"})
		return 0, false
	}
	return userID, true
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
