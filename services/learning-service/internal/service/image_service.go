package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/auth"
	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/storage"
)

// MaxImageBytes はアップロードできる画像の大きさ（5MB）
const MaxImageBytes = 5 << 20

// 受け付ける画像の種類と、保存するファイルの拡張子
var imageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type ImageService struct {
	storage storage.ImageStorage
	newName func() (string, error) // テストで差し替えるため
}

func NewImageService(s storage.ImageStorage) *ImageService {
	return &ImageService{storage: s, newName: randomName}
}

// Upload は画像を保存して表示用の URL を返す（管理者のみ）。
// 種類はファイル名や申告された Content-Type ではなく、中身の先頭のバイト列で判定する
func (s *ImageService) Upload(ctx context.Context, actor auth.Claims, data []byte) (model.UploadedImage, error) {
	if !actor.IsAdmin() {
		return model.UploadedImage{}, apperr.ErrForbidden
	}
	v := &validator{}
	if len(data) == 0 {
		v.add("画像ファイルを選んでください")
	} else if len(data) > MaxImageBytes {
		v.add("画像は5MB以下にしてください")
	}
	if err := v.err(); err != nil {
		return model.UploadedImage{}, err
	}
	contentType := http.DetectContentType(data)
	ext, ok := imageExtensions[contentType]
	if !ok {
		v.add("JPEG・PNG・WebP の画像を選んでください")
		return model.UploadedImage{}, v.err()
	}

	name, err := s.newName()
	if err != nil {
		return model.UploadedImage{}, fmt.Errorf("ファイル名の生成: %w", err)
	}
	url, err := s.storage.Save(ctx, name+ext, contentType, data)
	if err != nil {
		return model.UploadedImage{}, fmt.Errorf("画像の保存: %w", err)
	}
	return model.UploadedImage{URL: url, ContentType: contentType, Size: len(data)}, nil
}

// randomName は推測できない32文字のファイル名を作る（同じ名前の画像を上書きしないため）
func randomName() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
