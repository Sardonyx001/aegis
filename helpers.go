package aegis

import (
	"hash"

	"github.com/matthewhartstonge/argon2"
	"github.com/xdg-go/pbkdf2"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
)

// Salt size in bytes.
const saltSize = 16

func deriveKey(password string, salt []byte) []byte {
	// Default config takes care of salting for us
	argon := argon2.DefaultConfig()
	encodedPassword, _ := argon.HashEncoded([]byte(password))
	return pbkdf2.Key([]byte(encodedPassword), salt, 4096, chacha20poly1305.KeySize, func() hash.Hash {
		hashFunc, _ := blake2b.New512(nil)
		return hashFunc
	})
}
