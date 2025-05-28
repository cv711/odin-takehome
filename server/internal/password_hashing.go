package internal

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Combine salt and hash for storage
	combined := append(salt, hash...)
	return fmt.Sprintf("%x", combined), nil
}

func VerifyPassword(storedHash string, password string) bool {
	combined, err := hex.DecodeString(storedHash)
	if err != nil {
		return false
	}

	salt := combined[:32]
	storedHashBytes := combined[32:]

	// Generate new hash
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Compare hashes securely
	return subtle.ConstantTimeCompare(newHash, storedHashBytes) == 1
}
