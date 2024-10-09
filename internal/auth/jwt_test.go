package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tokenString == "" {
		t.Fatalf("expected a token string, got an empty string")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	// Create a valid JWT
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tokenString == "" {
		t.Fatalf("expected a token string, got an empty string")
	}

	// Validate the JWT
	validatedUserID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if validatedUserID != userID {
		t.Fatalf("expected userID %v, got %v", userID, validatedUserID)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "secret"
	invalidTokenString := "invalid.token.string"

	// Validate the invalid JWT
	_, err := ValidateJWT(invalidTokenString, tokenSecret)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := -time.Hour // Token already expired

	// Create an expired JWT
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tokenString == "" {
		t.Fatalf("expected a token string, got an empty string")
	}

	// Validate the expired JWT
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}
