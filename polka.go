package main

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (p *PgUserStore) GrantChirpyRed(ctx context.Context, userID uuid.UUID) error {
	err := p.db.AddChirpyRed(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &NotFoundError{Resource: "user", ID: userID.String()}
		}
		return &DatabaseError{Operation: "GrantChirpyRed", Err: err}
	}
	return nil
}
