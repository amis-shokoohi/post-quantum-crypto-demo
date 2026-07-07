package crypto

import (
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

const aesKeySize = 32

func DeriveKey(sharedSecret []byte) ([]byte, error) {
	kdf := hkdf.New(
		sha256.New,
		sharedSecret,
		nil,
		[]byte("chat-encryption"),
	)

	key := make([]byte, aesKeySize)

	if _, err := io.ReadFull(kdf, key); err != nil {
		return nil, fmt.Errorf("derive aes key: %w", err)
	}

	return key, nil
}
