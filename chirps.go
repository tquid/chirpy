package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tquid/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type ChirpStore interface {
	CreateChirp(ctx context.Context, body string, userID uuid.UUID) (*Chirp, error)
}

type PgChirpStore struct {
	db *database.Queries
}

func (p *PgChirpStore) CreateChirp(ctx context.Context, body string, userID uuid.UUID) (*Chirp, error) {
	params := database.CreateChirpParams{
		Body:   body,
		UserID: userID,
	}
	chirp, err := p.db.CreateChirp(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("unable to create chirp: %w", err)
	}
	return &Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}, nil
}
