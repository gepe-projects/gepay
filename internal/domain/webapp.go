package domain

import (
	"context"

	"github.com/google/uuid"
)

type WebappService interface {
	Signin(ctx context.Context, req SignInReq) (SigninRes, error)
	SignupLocal(ctx context.Context, req SignUpReq) error
	Me(ctx context.Context, id uuid.UUID) (UserWithIdentity, error)
}
