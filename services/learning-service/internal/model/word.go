package model

// Word は words テーブルの1行
type Word struct {
	ID               uint64  `json:"id"`
	English          string  `json:"english"`
	Japanese         string  `json:"japanese"`
	Example          *string `json:"example"`
	PhotoPath        *string `json:"photo_path"`
	IllustrationPath *string `json:"illustration_path"`
}

// WordInput は POST / PUT /words の body
type WordInput struct {
	English          string  `json:"english"`
	Japanese         string  `json:"japanese"`
	Example          *string `json:"example"`
	PhotoPath        *string `json:"photo_path"`
	IllustrationPath *string `json:"illustration_path"`
}
