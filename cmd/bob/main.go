package main

import (
	"log"
	"net"

	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/protocol"
	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/transport"
)

func main() {
	log.Println("[BOB] connecting to Alice")

	tcpConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer tcpConn.Close()

	conn := transport.NewConnection(tcpConn)

	cipher, err := protocol.ClientHandshake(conn)
	if err != nil {
		log.Fatalf("handshake failed: %v", err)
	}

	log.Println("[BOB] secure channel established")

	err = protocol.Chat(conn, cipher)
	if err != nil {
		log.Fatalf("chat ended: %v", err)
	}
}
