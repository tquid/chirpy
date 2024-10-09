package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{
		"Authorization": []string{"Bearer daljkhfladshfasdl"},
	}
	_, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetBearerToken_WrongHeaderLength(t *testing.T) {
	headers := http.Header{
		"Authorization": []string{"Bearer foo bar"},
	}
	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	headers = http.Header{
		"Authorization": []string{"Bearer"},
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetBearerToken_InvalidHeader(t *testing.T) {
	headers := http.Header{
		"Authorization": []string{"Quux adfasdfae"},
	}
	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
