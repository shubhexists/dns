package models

type DNSQuestion struct {
	QName  []byte
	QType  uint16
	QClass uint16
}
