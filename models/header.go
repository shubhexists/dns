package models

// I feel we are wasting a lot of bit space :(

type DNSHeader struct {
	PacketID uint16
	QR       uint8
	Opcode   uint8
	AA       uint8
	TC       uint8
	RD       uint8
	RA       uint8
	Z        uint8
	RCode    uint8
	QDCount  uint16
	ANCount  uint16
	NSCount  uint16
	ARCount  uint16
}
