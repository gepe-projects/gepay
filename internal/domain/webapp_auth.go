package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	SESSION_ID   = "session_id"
	SESSION_ROLE = "session_role"
)

type AuthService interface {
	Signin(ctx context.Context, req SignInReq) (SigninRes, error)
	SignupLocal(ctx context.Context, req SignUpReq) error
	Me(ctx context.Context, id uuid.UUID) (UserWithIdentity, error)
	Refresh(ctx context.Context, token string) (RefreshRes, error)
}

type SignInReq struct {
	Email    string `json:"email" form:"email"  binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
}

type SigninRes struct {
	UserWithIdentity UserWithIdentity `json:"user_with_identity"`
}
type RefreshRes struct {
	Refresh_token string `json:"refresh_token"`
}

type SignUpReq struct {
	Email        string `db:"email" json:"email" form:"email"  binding:"required,email"`
	Password     string `db:"password" json:"password" form:"password" binding:"required,min=8"`
	OwnerName    string `db:"owner_name" json:"owner_name" form:"owner_name" binding:"required"`
	BusinessName string `db:"business_name" json:"business_name" form:"business_name" binding:"required"`
	Role         string
	UserID       uuid.UUID
}

type TokenCLaims struct {
	Role   string
	Readme string
	Scopes []string
	jwt.RegisteredClaims
}

// type AuthRepository interface {
// 	FindUserByEmail(ctx context.Context, email string) (*User, error)
// 	SaveUser(ctx context.Context, user *User) error
// }
