package model

// Word は words テーブルの1行
type Word struct {
	ID               uint64  `json:"id"`
	English          string  `json:"english"`
	Japanese         string  `json:"japanese"`
	Example          *string `json:"example"`
	ExampleJapanese  *string `json:"example_japanese"`  // 例文の日本語訳
	PhotoPath        *string `json:"photo_path"`        // 画像の URL（POST /images が返したもの）
	IllustrationPath *string `json:"illustration_path"` // 同上
}

// WordInput は POST / PUT /words の body
type WordInput struct {
	English          string  `json:"english"`
	Japanese         string  `json:"japanese"`
	Example          *string `json:"example"`
	ExampleJapanese  *string `json:"example_japanese"`
	PhotoPath        *string `json:"photo_path"`
	IllustrationPath *string `json:"illustration_path"`
}
