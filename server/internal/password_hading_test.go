package internal_test

import (
	"testing"

	"github.com/cv711/odin-takehome/server/internal"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := internal.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == password {
		t.Fatalf("Hashed password should not be the same as the original password")
	}

	hashedPassword2, err := internal.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if hashedPassword == hashedPassword2 {
		t.Fatalf("Hashed passwords should not be the same for different calls")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := internal.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	valid := internal.VerifyPassword(hashedPassword, password)
	if !valid {
		t.Fatalf("Expected password verification to succeed")
	}

	invalid := internal.VerifyPassword(hashedPassword, "wrongpassword")
	if invalid {
		t.Fatalf("Expected password verification to fail")
	}
}
