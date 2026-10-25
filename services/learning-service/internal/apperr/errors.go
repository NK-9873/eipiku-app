// Package apperr は層をまたいで使うエラー。handler で HTTP ステータスに変換する
package apperr

import (
	"errors"
	"strings"
)

// ErrNotFound は指定した ID のデータが無いとき（404）
var ErrNotFound = errors.New("not found")

// ValidationError は入力が不正なとき（400）。Details に項目ごとの理由を入れる
type ValidationError struct {
	Details []string
}

func (e *ValidationError) Error() string {
	return "入力内容に誤りがあります: " + strings.Join(e.Details, "、")
}
