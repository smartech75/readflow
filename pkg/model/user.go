package model

import "time"

// User structure definition
type User struct {
	ID          *uint      `json:"id,omitempty"`
	Username    string     `json:"username,omitempty"`
	Enabled     bool       `json:"enabled,omitempty"`
	Plan        string     `json:"plan,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// UserForm structure definition
type UserForm struct {
	ID      uint
	Enabled *bool
	Plan    *string
}
