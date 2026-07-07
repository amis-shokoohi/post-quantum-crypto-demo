package main

import (
	"log"
	"net"

	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/protocol"
	"github.com/amis-shokoohi/post-quantum-crypto-demo/internal/transport"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	log.Println("[ALICE] listening on :8080")

	tcpConn, err := listener.Accept()
	if err != nil {
		log.Fatalf("accept connection: %v", err)
	}
	defer tcpConn.Close()

	log.Println("[ALICE] client connected")

	conn := transport.NewConnection(tcpConn)

	cipher, err := protocol.ServerHandshake(conn)
	if err != nil {
		log.Fatalf("handshake failed: %v", err)
	}

	log.Println("[ALICE] secure channel established")

	err = protocol.Chat(conn, cipher)
	if err != nil {
		log.Fatalf("chat ended: %v", err)
	}
}
