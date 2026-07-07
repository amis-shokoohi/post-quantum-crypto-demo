package transport

const (
	packetHeaderSize = 5
)

type Packet struct {
	Type uint8
	Data []byte
}
