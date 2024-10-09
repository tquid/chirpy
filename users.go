package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tquid/chirpy/internal/database"

	"github.com/google/uuid"
)

type UserStore interface {
	CreateUser(ctx context.Context, email string) (*User, error)
	DeleteAllUsers(ctx context.Context) (int64, error)
}

type PgUserStore struct {
	db *database.Queries
}

func (p *PgUserStore) CreateUser(ctx context.Context, email string) (*User, error) {
	user, err := p.db.CreateUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, nil
}

func (p *PgUserStore) DeleteAllUsers(ctx context.Context) (int64, error) {
	rows, err := p.db.DeleteAllUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to delete all users: %w", err)
	}
	return rows, nil
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
