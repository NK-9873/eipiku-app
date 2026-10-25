// Package auth は user-service が発行した JWT の検証と、ログイン中のユーザー情報の context での受け渡し。
// middleware がトークンを検証して context に入れ、handler が取り出す
package auth

import (
	"context"

	"eipiku/learning-service/internal/model"
)

// Claims はトークンから取り出したログイン中のユーザー
type Claims struct {
	UserID uint64
	Role   model.Role
}

func (c Claims) IsAdmin() bool {
	return c.Role == model.RoleAdmin
}

type claimsKey struct{}

func WithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, c)
}

// ClaimsFrom は context のユーザー情報を返す。入っていなければ ok=false
func ClaimsFrom(ctx context.Context) (c Claims, ok bool) {
	c, ok = ctx.Value(claimsKey{}).(Claims)
	return c, ok
}
