package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var ErrCipherTooShort = errors.New("ciphertext too short")

type Encryptor struct{}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}

func (e *Encryptor) Encrypt(data string, password string) (string, error) {
	hash := createPasswordHash(password)

	block, err := aes.NewCipher(hash)
	if err != nil {
		return "", fmt.Errorf("failed to create new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create new GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to read random data: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(data), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(data string, password string) (string, error) {
	hash := createPasswordHash(password)

	block, err := aes.NewCipher(hash)
	if err != nil {
		return "", fmt.Errorf("failed to create new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create new GCM: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrCipherTooShort
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

func createPasswordHash(password string) []byte {
	// Create a new hash.
	// TODO: Replace with Argon2id or Scrypt for better security against brute-force attacks.
	// Currently using SHA-256 for speed, but this is vulnerable to rainbow table attacks if the DB is leaked
	// (though mitigated by the fact that the server doesn't store the password, it only stores the data encrypted with it).
	// To fix this properly, we need to generate a random salt per bin, store it unencrypted, and use it in the KDF.
	h := sha256.New()
	// Write password to hash.
	h.Write([]byte(password))
	// Use the sum as the key.
	key := h.Sum(nil)
	// Return the key as a string.
	return key
}
