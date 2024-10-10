package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (p *PgUserStore) GetUserIDFromValidRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	refreshToken, err := p.db.GetRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, fmt.Errorf("refresh token not found in database")
		}
		return uuid.UUID{}, fmt.Errorf("error retrieving refresh token from database: %w", err)
	}
	if refreshToken.RevokedAt.Valid {
		return uuid.UUID{}, fmt.Errorf("refresh token is revoked")
	}
	return refreshToken.UserID, nil
}

func (p *PgUserStore) RevokeRefreshToken(ctx context.Context, token string) error {
	err := p.db.RevokeRefreshToken(ctx, token)
	if err != nil {
		return fmt.Errorf("refresh token revoke failed: %w", err)
	}
	return nil
}
