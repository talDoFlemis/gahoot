package ports

import (
	"context"
	"time"
)

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type UserIdentityInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AuthenticationManager interface {
	CreateUser(ctx context.Context, username, password string) (*UserIdentityInfo, error)
	AuthenticateUser(ctx context.Context, username, password string) (*UserIdentityInfo, error)
	DeleteUser(ctx context.Context, userId string) error
	UpdateUser(
		ctx context.Context,
		userId, username, password, email string,
	) (*UserIdentityInfo, error)
	CreateToken(ctx context.Context, userId string, expiresAt time.Time) (*TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
}

type LocalIDPUserEntity struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
}

type LocalIDPStorer interface {
	StoreUser(
		ctx context.Context,
		username, email, password string,
	) (*LocalIDPUserEntity, error)
	UpdateUser(
		ctx context.Context,
		userId, username, password, email string,
	) (*LocalIDPUserEntity, error)
	DeleteUser(ctx context.Context, userId string) error
	FindUserById(ctx context.Context, userId string) (*LocalIDPUserEntity, error)
}
