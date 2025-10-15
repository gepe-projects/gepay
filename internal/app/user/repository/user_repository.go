package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/ilhamgepe/gepay/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	log logger.Logger
	db  *sqlx.DB
}

func NewUserRepository(db *sqlx.DB, log logger.Logger) domain.UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

// CreateUserIdentityTx implements domain.UserRepository.
func (r *userRepository) CreateUserIdentityTx(ctx context.Context, tx *sqlx.Tx, req domain.SignUpReq) error {
	query := `
		INSERT INTO user_identities (user_id, provider, provider_id,email,password_hash)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.ExecContext(ctx, query, req.UserID, "local", req.Email, req.Email, req.Password)
	if err != nil {
		r.log.Errorf(err, "failed to create user identity")
		if utils.IsForeignKeyViolation(err) {
			return domain.ErrUserAlreadyExists
		}
		return domain.ErrInternalServerError
	}
	return nil
}

// CreateUserTx implements domain.UserRepository.
func (r *userRepository) CreateUserTx(ctx context.Context, tx *sqlx.Tx, req domain.SignUpReq) (uuid.UUID, error) {
	query := `
		INSERT INTO users (name,role)
		VALUES ($1, $2)
		RETURNING id
	`

	var id uuid.UUID
	err := tx.QueryRowContext(ctx, query, req.OwnerName, req.Role).Scan(&id)
	if err != nil {
		r.log.Errorf(err, "failed to create user")
		if utils.IsForeignKeyViolation(err) {
			return uuid.Nil, domain.ErrUserAlreadyExists
		}
		return uuid.Nil, domain.ErrInternalServerError
	}
	return id, nil
}

// FindUserByEmail implements domain.UserRepository.
func (r *userRepository) FindUserWithIdentityByEmail(ctx context.Context, email string) (domain.UserWithIdentity, error) {
	q := `
	SELECT 
		u.id AS "user.id",
		u.name AS "user.name",
		u.image_url AS "user.image_url",
		u.role AS "user.role",
		u.metadata AS "user.metadata",
		u.created_at AS "user.created_at",
		u.updated_at AS "user.updated_at",
		
		ui.id AS "useridentity.id",
		ui.user_id AS "useridentity.user_id",
		ui.provider AS "useridentity.provider",
		ui.provider_id AS "useridentity.provider_id",
		ui.email AS "useridentity.email",
		ui.phone AS "useridentity.phone",
		ui.password_hash AS "useridentity.password_hash",
		ui.verified AS "useridentity.verified",
		ui.last_login_at AS "useridentity.last_login_at",
		ui.created_at AS "useridentity.created_at",
		ui.updated_at AS "useridentity.updated_at"
	FROM users u
	JOIN user_identities ui ON ui.user_id = u.id
	WHERE ui.email = $1
	LIMIT 1
	`

	var uw domain.UserWithIdentity
	if err := r.db.GetContext(ctx, &uw, q, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uw, domain.ErrUserNotFound
		}
		r.log.Error(err, "error fetching user with identity")
		return uw, domain.ErrInternalServerError
	}

	return uw, nil
}

func (r *userRepository) FindUserWithIdentityByID(ctx context.Context, id uuid.UUID) (domain.UserWithIdentity, error) {
	q := `
	SELECT 
		u.id AS "user.id",
		u.name AS "user.name",
		u.image_url AS "user.image_url",
		u.role AS "user.role",
		u.metadata AS "user.metadata",
		u.created_at AS "user.created_at",
		u.updated_at AS "user.updated_at",
		
		ui.id AS "useridentity.id",
		ui.user_id AS "useridentity.user_id",
		ui.provider AS "useridentity.provider",
		ui.provider_id AS "useridentity.provider_id",
		ui.email AS "useridentity.email",
		ui.phone AS "useridentity.phone",
		ui.password_hash AS "useridentity.password_hash",
		ui.verified AS "useridentity.verified",
		ui.last_login_at AS "useridentity.last_login_at",
		ui.created_at AS "useridentity.created_at",
		ui.updated_at AS "useridentity.updated_at"
	FROM users u
	JOIN user_identities ui ON ui.user_id = u.id
	WHERE u.id = $1
	LIMIT 1
	`

	var uw domain.UserWithIdentity
	if err := r.db.GetContext(ctx, &uw, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uw, domain.ErrUserNotFound
		}
		r.log.Error(err, "error fetching user with identity")
		return uw, domain.ErrInternalServerError
	}

	return uw, nil
}
