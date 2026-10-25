// Package auth はリクエストしたユーザーの ID を context で受け渡す。
// ID を入れるのは middleware（今は X-User-ID ヘッダー、4回目で JWT の検証に置き換える）、
// 取り出すのは handler。置き換えても handler 側は変えずに済む
package auth

import "context"

type userIDKey struct{}

func WithUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserIDFrom は context のユーザー ID を返す。入っていなければ ok=false
func UserIDFrom(ctx context.Context) (userID uint64, ok bool) {
	userID, ok = ctx.Value(userIDKey{}).(uint64)
	return userID, ok
}
