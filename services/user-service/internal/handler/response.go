package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"eipiku/user-service/internal/apperr"
)

const maxBodyBytes = 1 << 20 // 1MB

type errorResponse struct {
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("レスポンスの書き込みに失敗: %v", err)
	}
}

func writeBadRequest(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusBadRequest, errorResponse{Error: message})
}

// writeError は service から返ったエラーを HTTP ステータスに変換する
func writeError(w http.ResponseWriter, err error) {
	var validationErr *apperr.ValidationError
	switch {
	case errors.As(err, &validationErr):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "入力内容に誤りがあります", Details: validationErr.Details})
	case errors.Is(err, apperr.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "見つかりません"})
	case errors.Is(err, apperr.ErrUnauthorized):
		// メールアドレスとパスワードのどちらが違うかは返さない
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "メールアドレスまたはパスワードが正しくありません"})
	case errors.Is(err, apperr.ErrForbidden):
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "この操作を行う権限がありません"})
	case errors.Is(err, apperr.ErrConflict):
		writeJSON(w, http.StatusConflict, errorResponse{Error: "このメールアドレスは既に登録されています"})
	default:
		log.Printf("サーバーエラー: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

// pathID は URL の {name} を ID として取り出す。数字でなければ 400 を書いて ok=false を返す
func pathID(w http.ResponseWriter, r *http.Request, name string) (id uint64, ok bool) {
	id, err := strconv.ParseUint(r.PathValue(name), 10, 64)
	if err != nil || id == 0 {
		writeBadRequest(w, name+" が不正です")
		return 0, false
	}
	return id, true
}

// decodeJSON は body を dst に読み込む。JSON として読めなければ 400 を書いて ok=false を返す
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) (ok bool) {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields() // 項目名の打ち間違いに気づけるようにする
	if err := dec.Decode(dst); err != nil {
		writeBadRequest(w, fmt.Sprintf("リクエストの JSON が不正です: %v", err))
		return false
	}
	return true
}
