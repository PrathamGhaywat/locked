package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	//Argon 2 params
	argon2Time = 3
	argon2Memory = 64 * 1024 //64 MB
	argon2Threads = 4
	argon2KeyLen = 32
	saltLen = 16
)

//dereives a 256-bit encryption key from a password using Argon2
func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)
}

//creates random salt for key derivation
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return  salt, nil
}

//converts salt bytes to hex string for storage
func SaltToString(salt []byte) string {
	return hex.EncodeToString(salt)
}

func SaltFromString(saltStr string) ([]byte, error) {
	salt, err := hex.DecodeString(saltStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}
	return salt, nil
}