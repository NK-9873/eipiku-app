package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"eipiku/user-service/internal/auth"
)

// TokenVerifier はトークンを検証して中身を返す（auth.JWT が満たす）
type TokenVerifier interface {
	Verify(token string) (auth.Claims, error)
}

// RequireAuth は Authorization: Bearer <JWT> を検証し、ユーザー情報を context に入れる。無い・不正・期限切れなら 401
func RequireAuth(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !found || token == "" {
				unauthorized(w)
				return
			}
			claims, err := verifier.Verify(token)
			if err != nil {
				unauthorized(w)
				return
			}
			next.ServeHTTP(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("WWW-Authenticate", `Bearer`)
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "ログインが必要です（トークンがないか、無効か、期限切れです）"})
}
