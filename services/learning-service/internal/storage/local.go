package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Local はディレクトリに保存し、baseURL + "/" + ファイル名 を URL として返す。
// 保存したファイルは GET /images/{name}（handler.ImageHandler）で表示する
type Local struct {
	dir     string
	baseURL string
}

func NewLocal(dir, baseURL string) (*Local, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("画像の保存先ディレクトリを作れません: %w", err)
	}
	return &Local{dir: dir, baseURL: strings.TrimSuffix(baseURL, "/")}, nil
}

func (l *Local) Save(ctx context.Context, name, contentType string, data []byte) (string, error) {
	// 書き込み途中のファイルが見えないよう、一時ファイルに書いてから名前を変える
	tmp, err := os.CreateTemp(l.dir, ".upload-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name()) // 名前を変えた後は何もしない
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(tmp.Name(), filepath.Join(l.dir, name)); err != nil {
		return "", err
	}
	return l.baseURL + "/" + name, nil
}

// Open は保存したファイルを開く。name は呼び出し側でファイル名の形式を確かめてから渡す
func (l *Local) Open(name string) (*os.File, error) {
	return os.Open(filepath.Join(l.dir, filepath.Base(name)))
}
