package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestLocal_SaveAndOpen(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "images") // 無いディレクトリでも作られる
	l, err := NewLocal(dir, "http://localhost:8082/images/")
	if err != nil {
		t.Fatal(err)
	}

	url, err := l.Save(context.Background(), "abc.png", "image/png", []byte("data"))
	if err != nil {
		t.Fatal(err)
	}
	if url != "http://localhost:8082/images/abc.png" {
		t.Errorf("url = %q", url)
	}

	f, err := l.Open("abc.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if b, _ := io.ReadAll(f); string(b) != "data" {
		t.Errorf("中身 = %q", b)
	}

	// 一時ファイルが残っていない
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("ディレクトリの中身 = %v", entries)
	}
}

func TestLocal_OpenOutsideDir(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "secret.txt"), []byte("secret"), 0o600)
	l, _ := NewLocal(filepath.Join(root, "images"), "http://localhost/images")

	if f, err := l.Open("../secret.txt"); err == nil {
		f.Close()
		t.Error("保存先の外のファイルが開けてしまった")
	}
}
