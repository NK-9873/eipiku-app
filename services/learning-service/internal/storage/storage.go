// Package storage はアップロードされた画像の保存先。
// 今はローカル（コンテナのボリューム）に保存しているが、ImageStorage を満たす S3 用の実装を作れば
// main.go の組み立てを変えるだけで切り替えられる
package storage

import "context"

// ImageStorage は画像を保存し、ブラウザで表示できる URL を返す
type ImageStorage interface {
	Save(ctx context.Context, name, contentType string, data []byte) (url string, err error)
}
