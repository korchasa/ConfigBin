package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const dayInSeconds = 60 * 60 * 24

func PasswordCookie(bid uuid.UUID, pass string) *http.Cookie {
	return &http.Cookie{
		Name:     bid.String(),
		Value:    pass,
		Path:     "/",
		MaxAge:   dayInSeconds,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func ReadPasswordCookie(r *http.Request, bid uuid.UUID) string {
	cookie, err := r.Cookie(bid.String())
	if err != nil {
		return ""
	}
	return cookie.Value
}

func PasswordFromHeader(r *http.Request) string {
	_, pass, _ := r.BasicAuth()
	return pass
}

// GeneratePassword generates a password of the given length.
func GeneratePassword(length int) (string, error) {
	upperCaseLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerCaseLetters := strings.ToLower(upperCaseLetters)
	digits := "0123456789"
	specials := "!@#$%&*()-_+=~"
	all := upperCaseLetters + lowerCaseLetters + digits + specials

	var password string
	for range length {
		number, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", fmt.Errorf("can't generate random letter: %w", err)
		}
		password += string(all[number.Int64()])
	}
	return password, nil
}
