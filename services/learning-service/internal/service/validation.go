package service

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"eipiku/learning-service/internal/apperr"
)

// validator は入力チェックの結果をためておき、最後にまとめて ValidationError にする
type validator struct {
	details []string
}

func (v *validator) add(format string, args ...any) {
	v.details = append(v.details, fmt.Sprintf(format, args...))
}

// required は前後の空白を除いた値を返す。空や文字数オーバーならエラーをためる
func (v *validator) required(field, value string, maxLen int) string {
	value = strings.TrimSpace(value)
	if value == "" {
		v.add("%s は必須です", field)
	} else if utf8.RuneCountInString(value) > maxLen {
		v.add("%s は%d文字以内で入力してください", field, maxLen)
	}
	return value
}

// optional は任意項目。未指定や空文字は nil（DB では NULL）にそろえる
func (v *validator) optional(field string, value *string, maxLen int) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	if utf8.RuneCountInString(trimmed) > maxLen {
		v.add("%s は%d文字以内で入力してください", field, maxLen)
	}
	return &trimmed
}

func (v *validator) err() error {
	if len(v.details) == 0 {
		return nil
	}
	return &apperr.ValidationError{Details: v.details}
}
