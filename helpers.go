package aegis

import (
	"hash"

	"github.com/xdg-go/pbkdf2"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
)

// Salt size in bytes.
const SALTSIZE = 16

func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, chacha20poly1305.KeySize, func() hash.Hash {
		hashFunc, _ := blake2b.New512(nil)
		return hashFunc
	})
}
