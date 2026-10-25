package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"eipiku/learning-service/internal/auth"
)

// UserFromHeader は X-User-ID ヘッダーのユーザー ID を context に入れる。無い・不正なら 401。
// 認証ができるまでの暫定で、4回目で JWT を検証する middleware に置き換える
func UserFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseUint(r.Header.Get("X-User-ID"), 10, 64)
		if err != nil || userID == 0 {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "ログインが必要です（X-User-ID ヘッダーがないか、値が不正です）"})
			return
		}
		next.ServeHTTP(w, r.WithContext(auth.WithUserID(r.Context(), userID)))
	})
}
