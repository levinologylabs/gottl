package hasher

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"strings"
)

const Prefix = "gottl_"

type Token struct {
	Raw  string
	Hash []byte
}

func generateToken(bits int) Token {
	randomBytes := make([]byte, bits)
	_, _ = rand.Read(randomBytes)

	plainText := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	plainText = Prefix + strings.ToLower(plainText)

	return Token{
		Raw:  plainText,
		Hash: HashToken(plainText),
	}
}

func GenerateShortToken() Token { return generateToken(16) }

// NewToken generates a new token prefixed with "rcp_"
func NewToken() Token { return generateToken(64) }

// HashToken hashes a token using SHA256
func HashToken(plainTextToken string) []byte {
	hash := sha256.Sum256([]byte(plainTextToken))
	return hash[:]
}
