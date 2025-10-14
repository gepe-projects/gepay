package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindUserWithIdentityByEmail(ctx context.Context, email string) (UserWithIdentity, error)
	FindUserWithIdentityByID(ctx context.Context, id uuid.UUID) (UserWithIdentity, error)
	CreateUserTx(ctx context.Context, tx *sqlx.Tx, req SignUpReq) (uuid.UUID, error)
	CreateUserIdentityTx(ctx context.Context, tx *sqlx.Tx, req SignUpReq) error
}

type UserService interface {
	FindUserWithIdentityByEmail(ctx context.Context, email string) (UserWithIdentity, error)
	FindUserWithIdentityByID(ctx context.Context, id uuid.UUID) (UserWithIdentity, error)
	CreateUserTx(ctx context.Context, tx *sqlx.Tx, req SignUpReq) (uuid.UUID, error)
	CreateUserIdentityTx(ctx context.Context, tx *sqlx.Tx, req SignUpReq) error
}

type UserWithIdentity struct {
	User         User
	UserIdentity UserIdentity
}

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	ImageURL  *string    `db:"image_url" json:"image_url,omitempty"`
	Role      string     `db:"role" json:"role"`
	Metadata  JSONB      `db:"metadata" json:"metadata"` // jsonb
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type UserIdentity struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	UserID       uuid.UUID  `db:"user_id" json:"user_id"`
	Provider     string     `db:"provider" json:"provider"`       // local, google, github, phone
	ProviderID   string     `db:"provider_id" json:"provider_id"` // unique ID dari provider
	Email        *string    `db:"email" json:"email,omitempty"`
	Phone        *string    `db:"phone" json:"phone,omitempty"`
	PasswordHash *string    `db:"password_hash" json:"-"`
	Verified     bool       `db:"verified" json:"verified"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
