package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"eipiku/learning-service/internal/model"
)

// 検証の方式は user-service（internal/auth/jwt.go）と同じにしている。片方を変えたらもう片方も合わせる

// MinSecretLength は JWT_SECRET に求める最低の長さ（HS256 の鍵として 256bit 以上）
const MinSecretLength = 32

// トークンの中身。発行するのは user-service（sub = ユーザー ID、role、exp = 24時間後）
type tokenClaims struct {
	Role model.Role `json:"role"`
	jwt.RegisteredClaims
}

// JWT は user-service が発行した HS256 のトークンを検証する（learning-service では発行しない）
type JWT struct {
	secret []byte
	now    func() time.Time // テストで時刻を差し替えるため
}

func NewJWT(secret string) (*JWT, error) {
	if len(secret) < MinSecretLength {
		return nil, fmt.Errorf("JWT_SECRET は%d文字以上にしてください（今は%d文字）", MinSecretLength, len(secret))
	}
	return &JWT{secret: []byte(secret), now: time.Now}, nil
}

// Verify は署名・アルゴリズム（HS256 のみ）・有効期限を確かめ、中身を返す
func (j *JWT) Verify(token string) (Claims, error) {
	var claims tokenClaims
	_, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (any, error) { return j.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(j.now),
	)
	if err != nil {
		return Claims{}, err
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || userID == 0 {
		return Claims{}, errors.New("sub が不正です")
	}
	if claims.Role != model.RoleUser && claims.Role != model.RoleAdmin {
		return Claims{}, errors.New("role が不正です")
	}
	return Claims{UserID: userID, Role: claims.Role}, nil
}
