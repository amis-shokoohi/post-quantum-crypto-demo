package protocol

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/crypto"
	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/transport"
)

func Chat(conn *transport.Connection, cipher *crypto.Cipher) error {
	done := make(chan struct{})

	var once sync.Once

	stop := func() {
		once.Do(func() {
			close(done)
		})
	}

	go sendLoop(conn, cipher, done, stop)

	recvLoop(conn, cipher, done, stop)

	return nil
}

func sendLoop(
	conn *transport.Connection,
	cipher *crypto.Cipher,
	done <-chan struct{},
	stop func(),
) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-done:
			return
		default:
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Printf("[CHAT] stdin error: %v", err)
			}
			stop()
			return
		}

		message := scanner.Bytes()

		encrypted, err := cipher.Encrypt(message)
		if err != nil {
			log.Printf("[CHAT] encrypt failed: %v", err)
			stop()
			return
		}

		err = conn.SendPacket(transport.Packet{
			Type: TypeEncryptedMessage,
			Data: encrypted,
		})
		if err != nil {
			log.Printf("[CHAT] send failed: %v", err)
			stop()
			return
		}
	}
}

func recvLoop(
	conn *transport.Connection,
	cipher *crypto.Cipher,
	done <-chan struct{},
	stop func(),
) {
	for {
		select {
		case <-done:
			return
		default:
		}

		packet, err := conn.RecvPacket()
		if err != nil {
			log.Printf("[CHAT] receive failed: %v", err)
			stop()
			return
		}

		if packet.Type != TypeEncryptedMessage {
			log.Printf("[CHAT] unexpected packet type: %d", packet.Type)
			continue
		}

		message, err := cipher.Decrypt(packet.Data)
		if err != nil {
			log.Printf("[CHAT] decrypt failed: %v", err)
			stop()
			return
		}

		fmt.Printf("peer: %s\n", string(message))
	}
}
