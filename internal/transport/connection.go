package transport

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) SendPacket(packet Packet) error {
	packetBytes := make([]byte, packetHeaderSize+len(packet.Data))

	packetBytes[0] = packet.Type

	binary.BigEndian.PutUint32(
		packetBytes[1:],
		uint32(len(packet.Data)),
	)

	copy(packetBytes[packetHeaderSize:], packet.Data)

	if _, err := c.conn.Write(packetBytes); err != nil {
		return fmt.Errorf("write packet: %w", err)
	}

	return nil
}

func (c *Connection) RecvPacket() (Packet, error) {
	header := make([]byte, packetHeaderSize)

	if _, err := io.ReadFull(c.conn, header); err != nil {
		return Packet{}, fmt.Errorf("read packet header: %w", err)
	}

	packetType := header[0]

	dataLen := binary.BigEndian.Uint32(header[1:])

	data := make([]byte, dataLen)

	if _, err := io.ReadFull(c.conn, data); err != nil {
		return Packet{}, fmt.Errorf("read packet data: %w", err)
	}

	return Packet{
		Type: packetType,
		Data: data,
	}, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
