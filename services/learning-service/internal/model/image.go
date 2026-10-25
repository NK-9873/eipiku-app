package model

// UploadedImage は POST /images のレスポンス。URL を words の photo_path / illustration_path に保存する
type UploadedImage struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
}
