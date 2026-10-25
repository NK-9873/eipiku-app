package handler

import (
	"fmt"
	"net/http"

	"eipiku/user-service/internal/auth"
	"eipiku/user-service/internal/model"
	"eipiku/user-service/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Register はルートを登録する。requireAuth は JWT を検証してユーザー情報を context に入れる middleware
func (h *UserHandler) Register(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler) {
	mux.HandleFunc("POST /signup", h.Signup)
	mux.HandleFunc("POST /login", h.Login)

	authed := func(pattern string, f http.HandlerFunc) {
		mux.Handle(pattern, requireAuth(f))
	}
	authed("GET /me", h.Me)
	authed("PUT /me/settings", h.UpdateSettings)
	// 管理者かどうかは service で確かめる
	authed("GET /users", h.List)
	authed("GET /users/{id}", h.Get)
}

// Signup は POST /signup
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var in model.SignupInput
	if !decodeJSON(w, r, &in) {
		return
	}
	user, err := h.service.Signup(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/users/%d", user.ID))
	writeJSON(w, http.StatusCreated, user)
}

// Login は POST /login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in model.LoginInput
	if !decodeJSON(w, r, &in) {
		return
	}
	result, err := h.service.Login(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// Me は GET /me
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	user, err := h.service.Me(r.Context(), actor)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// UpdateSettings は PUT /me/settings。body は {"image_type": "illustration"}
func (h *UserHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	var in model.SettingsInput
	if !decodeJSON(w, r, &in) {
		return
	}
	user, err := h.service.UpdateSettings(r.Context(), actor, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// List は GET /users（管理者のみ）
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	users, err := h.service.ListUsers(r.Context(), actor)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// Get は GET /users/{id}（管理者のみ）
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	id, ok := pathID(w, r, "id")
	if !ok {
		return
	}
	user, err := h.service.GetUser(r.Context(), actor, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// currentUser は middleware が context に入れたユーザー情報を取り出す
func currentUser(w http.ResponseWriter, r *http.Request) (auth.Claims, bool) {
	claims, ok := auth.ClaimsFrom(r.Context())
	if !ok {
		// middleware を通していないルートで呼ばれたとき（登録の書き忘れ）
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "ログインが必要です"})
		return auth.Claims{}, false
	}
	return claims, true
}
