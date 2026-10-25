package model

// Course は courses テーブルの1行
type Course struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"` // NULL のときは null
}

// CourseInput は POST / PUT /courses の body
type CourseInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// CourseWord はコースに含まれる単語と出題順（GET /courses/{id}/words の1件）
type CourseWord struct {
	DisplayOrder int `json:"display_order"`
	Word
}

// CourseWordsInput は PUT /courses/{id}/words の body。並び順がそのまま出題順になる
type CourseWordsInput struct {
	WordIDs []uint64 `json:"word_ids"`
}
