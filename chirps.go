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
	CreateChirp(context.Context, string, uuid.UUID) (*Chirp, error)
	GetChirps(context.Context) ([]*Chirp, error)
	GetChirpById(context.Context, uuid.UUID) (*Chirp, error)
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

func (p *PgChirpStore) GetChirps(ctx context.Context) ([]*Chirp, error) {
	DBchirps, err := p.db.GetChirps(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get chirps: %w", err)
	}
	var chirps []*Chirp
	for _, chirp := range DBchirps {
		chirps = append(chirps, &Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	return chirps, nil
}

func (p *PgChirpStore) GetChirpById(ctx context.Context, id uuid.UUID) (*Chirp, error) {
	DBChirp, err := p.db.GetChirpById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get chirp %s: %w", id, err)
	}
	return &Chirp{
		ID:        DBChirp.ID,
		CreatedAt: DBChirp.CreatedAt,
		UpdatedAt: DBChirp.UpdatedAt,
		Body:      DBChirp.Body,
		UserID:    DBChirp.UserID,
	}, nil
}
