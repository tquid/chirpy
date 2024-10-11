package main

import (
	"context"
	"database/sql"
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
	GetChirps(context.Context, SortDirection) ([]*Chirp, error)
	GetChirpsByUserId(context.Context, uuid.UUID, SortDirection) ([]*Chirp, error)
	GetChirpById(context.Context, uuid.UUID) (*Chirp, error)
	DeleteChirp(context.Context, uuid.UUID) error
}

type ChirpStoreError struct {
	Operation string
	Err       error
}

type PgChirpStore struct {
	db *database.Queries
}

func (cse *ChirpStoreError) Error() string {
	return fmt.Sprintf("%s: %v", cse.Operation, cse.Err)
}

func GetSortOrder(key string) string {
	if key == "desc" {
		return "DESC"
	}
	return "asc"
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

func (p *PgChirpStore) GetChirps(ctx context.Context, sort SortDirection) ([]*Chirp, error) {
	var dbFunc func(context.Context) ([]database.Chirp, error)
	if sort == SortDesc {
		dbFunc = p.db.GetChirpsDesc
	} else {
		dbFunc = p.db.GetChirps
	}
	DBchirps, err := dbFunc(ctx)
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

func (p *PgChirpStore) GetChirpsByUserId(ctx context.Context, userID uuid.UUID, sort SortDirection) ([]*Chirp, error) {
	var dbFunc func(context.Context, uuid.UUID) ([]database.Chirp, error)
	if sort == SortDesc {
		dbFunc = p.db.GetChirpsByUserIdDesc
	} else {
		dbFunc = p.db.GetChirpsByUserId
	}
	DBchirps, err := dbFunc(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Resource: "chirps by author", ID: userID.String()}
		}
		return nil, &DatabaseError{Operation: "GetChirpsByUserId", Err: err}
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
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Resource: "chirp", ID: id.String()}
		}
		return nil, &DatabaseError{Operation: "GetChirpById", Err: err}
	}
	return &Chirp{
		ID:        DBChirp.ID,
		CreatedAt: DBChirp.CreatedAt,
		UpdatedAt: DBChirp.UpdatedAt,
		Body:      DBChirp.Body,
		UserID:    DBChirp.UserID,
	}, nil
}

func (p *PgChirpStore) DeleteChirp(ctx context.Context, id uuid.UUID) error {
	err := p.db.DeleteChirp(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &NotFoundError{Resource: "chirp", ID: id.String()}
		}
		return &DatabaseError{Operation: "DeleteChirp", Err: err}
	}
	return nil
}
