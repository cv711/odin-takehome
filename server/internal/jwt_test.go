package internal_test

import (
	"testing"

	"github.com/cv711/odin-takehome/server/internal"
)

func TestGenerateJWTToken(t *testing.T) {
	// Test case 1: Valid user ID
	userId := "testUser"
	token, err := internal.GenerateJWTToken(userId)
	if err != nil {
		t.Fatalf("Failed to generate JWT token: %v", err)
	}

	// Validate the token
	claims, err := internal.ValidateJWTToken(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT token: %v", err)
	}

	if claims.Subject != userId {
		t.Fatalf("Expected subject to be %s, got %s", userId, claims.Subject)
	}
}
