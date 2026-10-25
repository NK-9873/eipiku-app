// Package apperr は層をまたいで使うエラー。handler で HTTP ステータスに変換する
package apperr

import (
	"errors"
	"strings"
)

var (
	// ErrNotFound は指定した ID のデータが無いとき（404）
	ErrNotFound = errors.New("not found")
	// ErrUnauthorized はログインに失敗したとき（401）
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden は権限が無い操作のとき（403）
	ErrForbidden = errors.New("forbidden")
	// ErrConflict は登録済みのデータと重複したとき（409）
	ErrConflict = errors.New("conflict")
)

// ValidationError は入力が不正なとき（400）。Details に項目ごとの理由を入れる
type ValidationError struct {
	Details []string
}

func (e *ValidationError) Error() string {
	return "入力内容に誤りがあります: " + strings.Join(e.Details, "、")
}
