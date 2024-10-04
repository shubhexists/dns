package models

type DNSQuestion struct {
	QName  []string
	QType  uint16
	QClass uint16
}
