package models

import (
	"time"
	"github.com/google/uuid"
)

type UserProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"not null" json:"email"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

// SignupRequest represents the signup payload
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

// AuthResponse represents the response from Supabase Auth
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	User         SupabaseUser `json:"user"`
}

type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// SigninRequest represents the signin payload
type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
