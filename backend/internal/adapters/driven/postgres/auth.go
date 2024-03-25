package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/taldoflemis/brain.test/internal/ports"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type LocalIDPPostgresStorer struct {
	pool *pgxpool.Pool
}

func NewLocalIDPPostgresStorer(pool *pgxpool.Pool) *LocalIDPPostgresStorer {
	return &LocalIDPPostgresStorer{
		pool: pool,
	}
}

func (s *LocalIDPPostgresStorer) StoreUser(
	ctx context.Context,
	username string,
	email string,
	password string,
) (*ports.LocalIDPUserEntity, error) {
	id := uuid.New()
	args := pgx.NamedArgs{
		"id":       id.String(),
		"username": username,
		"email":    email,
		"password": password,
	}
	insert := `INSERT INTO users (id, username, email, password) VALUES (@id, @username, @email, @password) RETURNING id`
	_, err := s.pool.Exec(ctx, insert, args)
	if err != nil {
		return nil, err
	}

	return &ports.LocalIDPUserEntity{
		ID:             id.String(),
		Username:       username,
		HashedPassword: password,
		Email:          email,
	}, nil
}

func (s *LocalIDPPostgresStorer) UpdateUser(
	ctx context.Context,
	userId string,
	username string,
	password string,
	email string,
) (*ports.LocalIDPUserEntity, error) {
	args := pgx.NamedArgs{
		"id":       userId,
		"username": username,
		"email":    email,
		"password": password,
	}
	updt := `UPDATE users SET username = @username, email = @email, password = @password WHERE id = @id RETURNING id`

	t, err := s.pool.Exec(ctx, updt, args)
	if err != nil {
		return nil, err
	}

	if t.RowsAffected() == 0 {
		return nil, ErrUserNotFound
	}

	return &ports.LocalIDPUserEntity{
		ID:             userId,
		Username:       username,
		Email:          email,
		HashedPassword: password,
	}, nil
}

func (s *LocalIDPPostgresStorer) DeleteUser(ctx context.Context, userId string) error {
	args := pgx.NamedArgs{
		"id": userId,
	}

	del := `DELETE FROM users WHERE id = @id`
	tag, err := s.pool.Exec(ctx, del, args)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *LocalIDPPostgresStorer) FindUserByUsername(
	ctx context.Context,
	username string,
) (*ports.LocalIDPUserEntity, error) {
	args := pgx.NamedArgs{
		"username": username,
	}

	query := `SELECT id, username, email, password FROM users WHERE username = @username`

	var user ports.LocalIDPUserEntity

	err := s.pool.QueryRow(ctx, query, args).
		Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
