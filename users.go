package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tquid/chirpy/internal/auth"
	"github.com/tquid/chirpy/internal/database"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var seconds int
	if err := json.Unmarshal(b, &seconds); err != nil {
		return fmt.Errorf("duration should be an integer representing seconds: %w", err)
	}
	d.Duration = time.Duration(seconds) * time.Second
	return nil
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStore interface {
	CreateUser(ctx context.Context, email string, password string) (*User, error)
	DeleteAllUsers(ctx context.Context) (int64, error)
	Login(ctx context.Context, jwtSecret string, loginParams LoginParams) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserIDFromValidRefreshToken(ctx context.Context, token string) (uuid.UUID, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	UpdateLoginCredentials(ctx context.Context, userID uuid.UUID, loginParams LoginParams) (*User, error)
	GrantChirpyRed(ctx context.Context, userID uuid.UUID) error
}

type PgUserStore struct {
	db *database.Queries
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (p *PgUserStore) CreateUser(ctx context.Context, email string, password string) (*User, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("can't create user: %w", err)
	}

	params := database.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPassword,
	}

	user, err := p.db.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}
	return &User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, nil
}

func (p *PgUserStore) DeleteAllUsers(ctx context.Context) (int64, error) {
	rows, err := p.db.DeleteAllUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to delete all users: %w", err)
	}
	return rows, nil
}

func (p *PgUserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := p.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("unable to get user: %w", err)
	}
	return &User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, nil
}

func (p *PgUserStore) Login(ctx context.Context, jwtSecret string, params LoginParams) (*User, error) {
	user, err := p.db.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &NotFoundError{Resource: "user", ID: params.Email}
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}
	refresh_token, err := auth.MakeRefreshToken()
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	}
	if err != nil {
		return nil, fmt.Errorf("can't make refresh token: %w", err)
	}
	err = p.db.CreateRefreshToken(ctx, refreshTokenParams)
	if err != nil {
		return nil, fmt.Errorf("adding refresh token to db failed: %w", err)
	}
	token, err := auth.MakeJWT(user.ID, jwtSecret, jwtExpireSeconds*time.Second)
	if err != nil {
		return nil, fmt.Errorf("can't make JWT token: %w", err)
	}
	return &User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refresh_token,
	}, nil
}

func (p *PgUserStore) UpdateLoginCredentials(ctx context.Context, userID uuid.UUID, params LoginParams) (*User, error) {
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}
	dbParams := database.UpdateLoginCredentialsParams{
		Email:          params.Email,
		ID:             userID,
		HashedPassword: hashedPassword,
	}

	user, err := p.db.UpdateLoginCredentials(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("updating db failed: %w", err)
	}
	return &User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, nil
}
