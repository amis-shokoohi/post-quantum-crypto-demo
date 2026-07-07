package crypto

import (
	"crypto/rand"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/mlkem/mlkem768"
)

func GenerateKeyPair() (kem.PublicKey, kem.PrivateKey, error) {
	return mlkem768.Scheme().GenerateKeyPair()
}

func Encapsulate(peerPubKey kem.PublicKey) (ciphertext []byte, sharedSecret []byte, err error) {
	seed := make([]byte, mlkem768.Scheme().EncapsulationSeedSize())
	if _, err := rand.Read(seed); err != nil {
		return nil, nil, err
	}
	ciphertext, sharedSecret, err = mlkem768.Scheme().EncapsulateDeterministically(peerPubKey, seed)
	return ciphertext, sharedSecret, err
}

func Decapsulate(privateKey kem.PrivateKey, ciphertext []byte) (sharedSecret []byte, err error) {
	return mlkem768.Scheme().Decapsulate(privateKey, ciphertext)
}

func UnmarshalPublicKey(data []byte) (kem.PublicKey, error) {
	return mlkem768.Scheme().UnmarshalBinaryPublicKey(data)
}
