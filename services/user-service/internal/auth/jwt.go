package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"eipiku/user-service/internal/model"
)

// MinSecretLength は JWT_SECRET に求める最低の長さ（HS256 の鍵として 256bit 以上）
const MinSecretLength = 32

const tokenTTL = 24 * time.Hour

type tokenClaims struct {
	Role model.Role `json:"role"`
	jwt.RegisteredClaims
}

// JWT は HS256 のトークンを発行・検証する
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

// Issue はユーザー ID（sub）・role・有効期限（24時間）を入れたトークンを発行する
func (j *JWT) Issue(userID uint64, role model.Role) (string, time.Time, error) {
	now := j.now()
	expiresAt := now.Add(tokenTTL)
	claims := tokenClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
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
