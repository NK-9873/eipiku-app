package model

// Course は courses テーブルの1行
type Course struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"` // NULL のときは null
}
