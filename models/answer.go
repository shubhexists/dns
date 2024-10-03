package models

type DNSAnswer struct {
	Name     []string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []int
}
