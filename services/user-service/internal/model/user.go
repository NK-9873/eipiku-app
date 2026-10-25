package model

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (r Role) Valid() bool {
	return r == RoleUser || r == RoleAdmin
}

// ImageType は単語学習で表示する画像（写真・イラスト）
type ImageType string

const (
	ImageTypePhoto        ImageType = "photo"
	ImageTypeIllustration ImageType = "illustration"
)

func (t ImageType) Valid() bool {
	return t == ImageTypePhoto || t == ImageTypeIllustration
}

// User は users テーブルの1行。PasswordHash は JSON に出さない
type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	ImageType    ImageType `json:"image_type"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser は登録時に repository に渡す値（パスワードはハッシュ化済み）
type NewUser struct {
	Name         string
	Email        string
	PasswordHash string
	Role         Role
}

// UserUpdate は管理者による更新で repository に渡す値。PasswordHash が nil ならパスワードは変えない
type UserUpdate struct {
	Name         string
	Email        string
	Role         Role
	PasswordHash *string
}

// AdminUserInput は POST /users・PUT /users/{id} の body（管理者用）。
// 作成では role を省略すると一般ユーザー。更新では password が空ならパスワードを変えない
type AdminUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

// SignupInput は POST /signup の body
type SignupInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginInput は POST /login の body
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResult は POST /login のレスポンス
type LoginResult struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"` // Authorization ヘッダーに付ける種類（Bearer）
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

// SettingsInput は PUT /me/settings の body
type SettingsInput struct {
	ImageType ImageType `json:"image_type"`
}
