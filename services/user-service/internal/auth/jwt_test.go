package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"eipiku/user-service/internal/model"
)

const testSecret = "test-secret-0123456789-0123456789-abcdef"

func newTestJWT(t *testing.T, now time.Time) *JWT {
	t.Helper()
	j, err := NewJWT(testSecret)
	if err != nil {
		t.Fatal(err)
	}
	j.now = func() time.Time { return now }
	return j
}

func TestJWT_IssueAndVerify(t *testing.T) {
	now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
	j := newTestJWT(t, now)

	token, expiresAt, err := j.Issue(7, model.RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if !expiresAt.Equal(now.Add(24 * time.Hour)) {
		t.Errorf("有効期限は24時間後のはずが %v", expiresAt)
	}

	claims, err := j.Verify(token)
	if err != nil {
		t.Fatalf("発行したトークンを検証できない: %v", err)
	}
	if claims.UserID != 7 || claims.Role != model.RoleAdmin {
		t.Errorf("claims = %#v", claims)
	}

	t.Run("24時間を過ぎたら無効", func(t *testing.T) {
		j.now = func() time.Time { return now.Add(24*time.Hour + time.Second) }
		defer func() { j.now = func() time.Time { return now } }()
		if _, err := j.Verify(token); err == nil {
			t.Error("期限切れなのに検証が通った")
		}
	})

	t.Run("別の鍵で署名したトークンは無効", func(t *testing.T) {
		other, _ := NewJWT(strings.Repeat("x", MinSecretLength))
		other.now = j.now
		forged, _, _ := other.Issue(7, model.RoleAdmin)
		if _, err := j.Verify(forged); err == nil {
			t.Error("別の鍵のトークンが通った")
		}
	})

	t.Run("HS256 以外（alg=none）は無効", func(t *testing.T) {
		none, _ := jwt.NewWithClaims(jwt.SigningMethodNone, tokenClaims{
			Role:             model.RoleAdmin,
			RegisteredClaims: jwt.RegisteredClaims{Subject: "7", ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour))},
		}).SignedString(jwt.UnsafeAllowNoneSignatureType)
		if _, err := j.Verify(none); err == nil {
			t.Error("alg=none のトークンが通った")
		}
	})
}

func TestNewJWT_ShortSecret(t *testing.T) {
	if _, err := NewJWT("change_me"); err == nil {
		t.Error("短い鍵が受け付けられた")
	}
}
