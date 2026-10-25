package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"eipiku/learning-service/internal/auth"
)

const testSecret = "test-secret-0123456789-0123456789-abcdef"

// sign は user-service と同じ形のトークンを作る（sub・role・exp）
func sign(t *testing.T, method jwt.SigningMethod, key any, sub, role string, exp time.Time) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, jwt.MapClaims{
		"sub":  sub,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  exp.Unix(),
	}).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestRequireAuth(t *testing.T) {
	verifier, err := auth.NewJWT(testSecret)
	if err != nil {
		t.Fatal(err)
	}
	valid := sign(t, jwt.SigningMethodHS256, []byte(testSecret), "7", "admin", time.Now().Add(time.Hour))

	// 後ろの handler では、context に入ったユーザー情報を記録する
	var gotClaims auth.Claims
	var reached bool
	h := RequireAuth(verifier)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		gotClaims, _ = auth.ClaimsFrom(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	tests := []struct {
		name          string
		authorization string
		wantStatus    int
	}{
		{"正しいトークン", "Bearer " + valid, http.StatusNoContent},
		{"トークンなし", "", http.StatusUnauthorized},
		{"Bearer が付いていない", valid, http.StatusUnauthorized},
		{"改ざん（中身の role を書き換え）", "Bearer " + tamperRole(t, valid), http.StatusUnauthorized},
		{"改ざん（署名を1文字変更）", "Bearer " + tamperSignature(valid), http.StatusUnauthorized},
		{"期限切れ", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testSecret), "7", "admin", time.Now().Add(-time.Minute)), http.StatusUnauthorized},
		{"別の鍵で署名", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(strings.Repeat("x", 40)), "7", "admin", time.Now().Add(time.Hour)), http.StatusUnauthorized},
		{"alg=none", "Bearer " + sign(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, "7", "admin", time.Now().Add(time.Hour)), http.StatusUnauthorized},
		{"HS256 以外のアルゴリズム（HS512）", "Bearer " + sign(t, jwt.SigningMethodHS512, []byte(testSecret), "7", "admin", time.Now().Add(time.Hour)), http.StatusUnauthorized},
		{"role が不正", "Bearer " + sign(t, jwt.SigningMethodHS256, []byte(testSecret), "7", "root", time.Now().Add(time.Hour)), http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reached, gotClaims = false, auth.Claims{}
			req := httptest.NewRequest(http.MethodGet, "/me/summary", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d（body: %s）", rec.Code, tt.wantStatus, rec.Body)
			}
			if tt.wantStatus == http.StatusUnauthorized {
				if reached {
					t.Error("401 なのに後ろの handler が呼ばれた")
				}
				return
			}
			if gotClaims.UserID != 7 || !gotClaims.IsAdmin() {
				t.Errorf("context のユーザー情報 = %#v", gotClaims)
			}
		})
	}
}

// tamperRole はトークンの中身を role=admin・sub=999 に書き換える（署名はそのまま）
func tamperRole(t *testing.T, token string) string {
	t.Helper()
	parts := strings.Split(token, ".")
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"999","role":"admin","exp":4102444800}`))
	return parts[0] + "." + payload + "." + parts[2]
}

// tamperSignature は署名の先頭の1文字を変える。
// 末尾の文字はパディングのビットを含み、変えても署名が同じになることがあるので先頭を使う
func tamperSignature(token string) string {
	parts := strings.Split(token, ".")
	sig := []byte(parts[2])
	if sig[0] == 'A' {
		sig[0] = 'B'
	} else {
		sig[0] = 'A'
	}
	return parts[0] + "." + parts[1] + "." + string(sig)
}
