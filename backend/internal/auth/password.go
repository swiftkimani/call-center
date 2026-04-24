package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 2
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 16
)

var ErrInvalidHash = errors.New("invalid password hash format")

func HashPassword(plain string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}
	hash := argon2.IDKey([]byte(plain), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

func VerifyPassword(encodedHash, plain string) (bool, error) {
	salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	candidate := argon2.IDKey([]byte(plain), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return subtle.ConstantTimeCompare(hash, candidate) == 1, nil
}

func decodeHash(encoded string) (salt, hash []byte, err error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return nil, nil, ErrInvalidHash
	}
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, ErrInvalidHash
	}
	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, ErrInvalidHash
	}
	return salt, hash, nil
}
