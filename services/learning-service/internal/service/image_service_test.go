package service

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"eipiku/learning-service/internal/apperr"
)

// mockImageStorage は storage.ImageStorage のモック。保存した内容を記録する
type mockImageStorage struct {
	saved       bool
	name        string
	contentType string
}

func (m *mockImageStorage) Save(ctx context.Context, name, contentType string, data []byte) (string, error) {
	m.saved, m.name, m.contentType = true, name, contentType
	return "http://localhost:8082/images/" + name, nil
}

// 各形式のファイルの先頭（中身の判定に使うバイト列）
var (
	pngHeader  = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")
	jpegHeader = []byte("\xff\xd8\xff\xe0\x00\x10JFIF\x00")
	webpHeader = []byte("RIFF\x24\x00\x00\x00WEBPVP8 ")
)

func newTestImageService(s *mockImageStorage) *ImageService {
	svc := NewImageService(s)
	svc.newName = func() (string, error) { return "0123456789abcdef0123456789abcdef", nil }
	return svc
}

func TestImageService_Upload(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantType string
		wantFile string
	}{
		{"PNG", pngHeader, "image/png", "0123456789abcdef0123456789abcdef.png"},
		{"JPEG", jpegHeader, "image/jpeg", "0123456789abcdef0123456789abcdef.jpg"},
		{"WebP", webpHeader, "image/webp", "0123456789abcdef0123456789abcdef.webp"},
	}
	for _, tt := range tests {
		t.Run(tt.name+"は中身から判定した拡張子で保存し、URL を返す", func(t *testing.T) {
			s := &mockImageStorage{}
			got, err := newTestImageService(s).Upload(context.Background(), adminActor, tt.data)
			if err != nil {
				t.Fatalf("エラーが返った: %v", err)
			}
			if s.name != tt.wantFile || s.contentType != tt.wantType {
				t.Errorf("保存した名前・種類 = %q %q", s.name, s.contentType)
			}
			if got.URL != "http://localhost:8082/images/"+tt.wantFile || got.ContentType != tt.wantType || got.Size != len(tt.data) {
				t.Errorf("got %#v", got)
			}
		})
	}

	rejects := []struct {
		name       string
		data       []byte
		wantDetail string
	}{
		{"空のファイル", nil, "画像ファイルを選んでください"},
		{"5MBを超える", append(append([]byte{}, pngHeader...), bytes.Repeat([]byte{0}, MaxImageBytes)...), "5MB以下"},
		{"テキスト", []byte("hello, this is not an image"), "JPEG・PNG・WebP"},
		{"GIF", []byte("GIF89a\x01\x00\x01\x00"), "JPEG・PNG・WebP"},
		{"SVG（拡張子を偽っても中身で判定）", []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`), "JPEG・PNG・WebP"},
	}
	for _, tt := range rejects {
		t.Run(tt.name+"は ValidationError で、保存しない", func(t *testing.T) {
			s := &mockImageStorage{}
			_, err := newTestImageService(s).Upload(context.Background(), adminActor, tt.data)
			assertValidationError(t, err, tt.wantDetail)
			if s.saved {
				t.Error("保存されてしまった")
			}
		})
	}

	t.Run("ちょうど5MBは受け付ける", func(t *testing.T) {
		data := append(append([]byte{}, pngHeader...), bytes.Repeat([]byte{0}, MaxImageBytes-len(pngHeader))...)
		if _, err := newTestImageService(&mockImageStorage{}).Upload(context.Background(), adminActor, data); err != nil {
			t.Errorf("エラーが返った: %v", err)
		}
	})

	t.Run("一般ユーザーは ErrForbidden で、保存しない", func(t *testing.T) {
		s := &mockImageStorage{}
		_, err := newTestImageService(s).Upload(context.Background(), userActor, pngHeader)
		if !errors.Is(err, apperr.ErrForbidden) || s.saved {
			t.Errorf("err %v, saved %v", err, s.saved)
		}
	})
}

func TestRandomName(t *testing.T) {
	a, err1 := randomName()
	b, err2 := randomName()
	if err1 != nil || err2 != nil || len(a) != 32 || a == b {
		t.Errorf("a=%q b=%q", a, b)
	}
}
