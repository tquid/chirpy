package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("api key not found")
	}
	headerParts := strings.Fields(authHeader)
	if len(headerParts) != 2 {
		return "", fmt.Errorf("wrong number of fields in Authorization header: got %d, expected 2", len(headerParts))
	}
	if headerParts[0] != "ApiKey" {
		return "", fmt.Errorf("invalid Authorization header '%s'", headerParts[0])
	}
	return headerParts[1], nil
}
