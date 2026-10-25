// Package auth は JWT の発行・検証と、ログイン中のユーザー情報の context での受け渡し
package auth

import (
	"context"

	"eipiku/user-service/internal/model"
)

// Claims はトークンから取り出したログイン中のユーザー
type Claims struct {
	UserID uint64
	Role   model.Role
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
