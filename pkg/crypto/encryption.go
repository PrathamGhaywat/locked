package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	nonceLenGCM = 12 //96-bit nonce for GCM
)

// encrypts plaintext using AES-256-GCM
// returns: nonce + ciphertext + tag (concatenated)
func EncryptData(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes for AES-256, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceLenGCM)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	//returns ciphertext + tag
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypts ciphertext using AES-256-GCM
// expects: nonce + ciphertext + tag (concentated)
func DecryptData(key []byte, ciphertext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes for AES-256, got: %d", len(key))
	}

	if len(ciphertext) < nonceLenGCM {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := ciphertext[:nonceLenGCM]
	encryptedData := ciphertext[nonceLenGCM:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong password?): %w", err)
	}

	return plaintext, nil
}

// encrypts data from reader and writes to writer
func EncryptStream(key []byte, reader io.Reader, writer io.Writer) error {
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceLenGCM)
	_, err = rand.Read(nonce)
	if err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	//write nonce first.
	_, err = writer.Write(nonce)
	if err != nil {
		return fmt.Errorf("failed to write nonce: %w", err)
	}

	//read all data and encrypt as single block for security.
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	encrypted := gcm.Seal(nil, nonce, data, nil)
	_, err = writer.Write(encrypted)
	if err != nil {
		return fmt.Errorf("failed to write encrypted data: %w", err)
	}

	return nil
}

// decrypts data from reader and writes to writer
func DecryptStream(key []byte, reader io.Reader, writer io.Writer) error {
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceLenGCM)
	_, err = reader.Read(nonce)
	if err != nil {
		return fmt.Errorf("failed to read nonce: %w", err)
	}

	//read all encrypted data.
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read encrypted data: %w", err)
	}

	decrypted, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return fmt.Errorf("decryption failed (wrong password?): %w", err)
	}

	_, err = writer.Write(decrypted)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data: %w", err)
	}

	return nil
}
