package models

type DNSAnswer struct {
	Name     []byte
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}
