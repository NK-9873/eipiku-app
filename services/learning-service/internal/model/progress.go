package model

import "time"

// Mastery は word_progresses.mastery（苦手・微妙・覚えた）
type Mastery string

const (
	MasteryWeak     Mastery = "weak"
	MasteryUnsure   Mastery = "unsure"
	MasteryMastered Mastery = "mastered"
)

func (m Mastery) Valid() bool {
	return m == MasteryWeak || m == MasteryUnsure || m == MasteryMastered
}

// WordProgress は word_progresses テーブルの1行
type WordProgress struct {
	UserID    uint64    `json:"user_id"`
	WordID    uint64    `json:"word_id"`
	Mastery   Mastery   `json:"mastery"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WordProgressInput は PUT /me/word-progresses/{wordId} の body
type WordProgressInput struct {
	Mastery Mastery `json:"mastery"`
}

// ReviewWord は復習対象の単語（単語の情報と今の習得度）
type ReviewWord struct {
	Word
	Mastery   Mastery   `json:"mastery"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProgressFilter は習得度の絞り込み条件。CourseID が nil なら全コース、Masteries が空なら全習得度
type ProgressFilter struct {
	CourseID  *uint64
	Masteries []Mastery
}

// CourseProgress は course_progresses テーブルの1行。
// CurrentPosition は「何語目まで学習したか」（0 なら未学習、単語数と同じなら最後まで終わっている）
type CourseProgress struct {
	UserID          uint64    `json:"user_id"`
	CourseID        uint64    `json:"course_id"`
	CurrentPosition int       `json:"current_position"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CourseProgressInput は PUT /me/course-progresses/{courseId} の body。未指定を見分けるためポインタにしている
type CourseProgressInput struct {
	CurrentPosition *int `json:"current_position"`
}

// コースの進み具合の状態
const (
	CourseNotStarted = "not_started"
	CourseInProgress = "in_progress"
	CourseCompleted  = "completed"
)

// Summary はマイページ用の学習状況（GET /me/summary）
type Summary struct {
	LearnedWords int             `json:"learned_words"` // 1回でも評価した単語の数
	TotalWords   int             `json:"total_words"`
	Weak         int             `json:"weak"`
	Unsure       int             `json:"unsure"`
	Mastered     int             `json:"mastered"`
	Courses      []CourseSummary `json:"courses"`
}

type CourseSummary struct {
	CourseID        uint64     `json:"course_id"`
	Name            string     `json:"name"`
	TotalWords      int        `json:"total_words"`
	CurrentPosition int        `json:"current_position"` // 単語数を超えないように丸めた値
	Status          string     `json:"status"`           // not_started / in_progress / completed
	UpdatedAt       *time.Time `json:"updated_at"`       // 未学習なら null
}
