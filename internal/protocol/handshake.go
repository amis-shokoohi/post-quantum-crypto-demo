package protocol

import (
	"fmt"
	"log"

	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/crypto"
	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/transport"
)

func ClientHandshake(conn *transport.Connection) (*crypto.Cipher, error) {
	log.Println("[BOB] starting handshake")

	// Bob generates ML-KEM keypair
	log.Println("[BOB] generating ML-KEM-768 keypair")
	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generate keypair: %w", err)
	}

	// Send public key to Alice
	pubBytes, err := pub.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	log.Println("[BOB] sending ML-KEM public key to Alice")
	err = conn.SendPacket(transport.Packet{
		Type: TypePublicKey,
		Data: pubBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("send public key: %w", err)
	}

	// Receive ciphertext from Alice
	log.Println("[BOB] waiting for KEM ciphertext")
	packet, err := conn.RecvPacket()
	if err != nil {
		return nil, fmt.Errorf("receive ciphertext: %w", err)
	}

	if packet.Type != TypeKEMCiphertext {
		return nil, fmt.Errorf("unexpected packet type: %d", packet.Type)
	}

	log.Println("[BOB] received KEM ciphertext from Alice")

	// Derive shared secret
	log.Println("[BOB] decapsulating shared secret")
	sharedSecret, err := crypto.Decapsulate(priv, packet.Data)
	if err != nil {
		return nil, fmt.Errorf("decapsulate: %w", err)
	}

	// Derive AES key
	log.Println("[BOB] deriving AES key")
	key, err := crypto.DeriveKey(sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("derive aes key: %w", err)
	}

	// Create AES-GCM session
	log.Println("[BOB] creating AES-GCM cipher session")
	cipher, err := crypto.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	log.Println("[BOB] handshake complete")

	return cipher, nil
}

func ServerHandshake(conn *transport.Connection) (*crypto.Cipher, error) {
	log.Println("[ALICE] starting handshake")

	// Receive Bob's public key
	log.Println("[ALICE] waiting for ML-KEM public key")
	packet, err := conn.RecvPacket()
	if err != nil {
		return nil, fmt.Errorf("receive public key: %w", err)
	}

	if packet.Type != TypePublicKey {
		return nil, fmt.Errorf("unexpected packet type: %d", packet.Type)
	}

	log.Println("[ALICE] received ML-KEM public key from Bob")

	// Parse public key
	log.Println("[ALICE] parsing public key")
	pub, err := crypto.UnmarshalPublicKey(packet.Data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal public key: %w", err)
	}

	// Encapsulate shared secret
	log.Println("[ALICE] encapsulating shared secret")
	ciphertext, sharedSecret, err := crypto.Encapsulate(pub)
	if err != nil {
		return nil, fmt.Errorf("encapsulate: %w", err)
	}

	// Send ciphertext to Bob
	log.Println("[ALICE] sending KEM ciphertext to Bob")
	err = conn.SendPacket(transport.Packet{
		Type: TypeKEMCiphertext,
		Data: ciphertext,
	})
	if err != nil {
		return nil, fmt.Errorf("send ciphertext: %w", err)
	}

	// Derive AES key
	log.Println("[ALICE] deriving AES key")
	key, err := crypto.DeriveKey(sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("derive aes key: %w", err)
	}

	// Create AES-GCM session
	log.Println("[ALICE] creating AES-GCM cipher session")
	cipher, err := crypto.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	log.Println("[ALICE] handshake complete")

	return cipher, nil
}
