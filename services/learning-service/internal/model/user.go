package model

// Role は user-service の users.role（トークンに入っている）
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)
