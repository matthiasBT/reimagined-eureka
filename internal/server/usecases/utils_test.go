package usecases

import (
	"encoding/base64"
	"testing"
)

func TestGenerateSessionToken(t *testing.T) {
	token := generateSessionToken()

	if token == "" {
		t.Errorf("Generated token should not be empty")
	}

	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("Failed to decode base64 token: %v", err)
	}

	if len(data) != 32 {
		t.Errorf("Expected decoded token to be 32 bytes, got %d", len(data))
	}
}
