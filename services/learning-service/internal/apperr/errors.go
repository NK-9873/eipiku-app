// Package apperr は層をまたいで使うエラー。handler で HTTP ステータスに変換する
package apperr

import (
	"errors"
	"strings"
)

var (
	// ErrNotFound は指定した ID のデータが無いとき（404）
	ErrNotFound = errors.New("not found")
	// ErrForbidden は権限が無い操作のとき（403）
	ErrForbidden = errors.New("forbidden")
)

// ValidationError は入力が不正なとき（400）。Details に項目ごとの理由を入れる
type ValidationError struct {
	Details []string
}

func (e *ValidationError) Error() string {
	return "入力内容に誤りがあります: " + strings.Join(e.Details, "、")
}
