package domain

import (
	"context"

	"github.com/iqdf/benjerry-service/common/auth"
)

// User domain
type User struct {
	Username       string
	HashPassword   string
	Authorizations []auth.Authorization
}

// UserService ...
type UserService interface {
	RegisterUser(ctx context.Context, username, hashpass string, isAdmin bool) error
	LoginUser(ctx context.Context, username, hashpass string) (User, error)
}

// UserRepository ...
type UserRepository interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, username string) (User, error)
}
