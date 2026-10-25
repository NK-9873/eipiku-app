package handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"eipiku/learning-service/internal/service"
)

// ImageOpener は保存した画像を開く（storage.Local が満たす）
type ImageOpener interface {
	Open(name string) (*os.File, error)
}

type ImageHandler struct {
	service *service.ImageService
	opener  ImageOpener // ローカル保存のときだけ使う。S3 に切り替えたら nil にして GET /images を登録しない
}

func NewImageHandler(s *service.ImageService, opener ImageOpener) *ImageHandler {
	return &ImageHandler{service: s, opener: opener}
}

// Register はルートを登録する。アップロードは管理者のみ（service で確かめる）。
// 表示は <img> タグから読み込むので、トークンなしで見られるようにしている
func (h *ImageHandler) Register(mux *http.ServeMux, requireAuth func(http.Handler) http.Handler) {
	mux.Handle("POST /images", requireAuth(http.HandlerFunc(h.Upload)))
	if h.opener != nil {
		mux.HandleFunc("GET /images/{name}", h.Serve)
	}
}

// Upload は POST /images。multipart/form-data の file 項目で画像を受け取る
func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	// 画像の上限＋multipart の区切りなどのぶんだけ読む
	r.Body = http.MaxBytesReader(w, r.Body, service.MaxImageBytes+1<<20)
	file, _, err := r.FormFile("file")
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeBadRequest(w, "画像は5MB以下にしてください")
			return
		}
		writeBadRequest(w, "画像ファイルを file という項目名で送ってください（multipart/form-data）")
		return
	}
	defer file.Close()
	// 上限を1バイトでも超えたら service で 400 にする
	data, err := io.ReadAll(io.LimitReader(file, service.MaxImageBytes+1))
	if err != nil {
		writeBadRequest(w, "画像ファイルを読み込めませんでした")
		return
	}

	image, err := h.service.Upload(r.Context(), actor, data)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, image)
}

// 保存時に付けた名前（32文字の16進数＋拡張子）だけを受け付け、ほかのファイルは読ませない
var imageNamePattern = regexp.MustCompile(`^[0-9a-f]{32}\.(jpg|png|webp)$`)

var imageContentTypes = map[string]string{".jpg": "image/jpeg", ".png": "image/png", ".webp": "image/webp"}

// Serve は GET /images/{name}
func (h *ImageHandler) Serve(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if !imageNamePattern.MatchString(name) {
		http.NotFound(w, r)
		return
	}
	f, err := h.opener.Open(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", imageContentTypes[filepath.Ext(name)])
	w.Header().Set("X-Content-Type-Options", "nosniff")
	// 同じ名前のファイルは作り直さないので、長くキャッシュしてよい
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeContent(w, r, name, info.ModTime(), f)
}
