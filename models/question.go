package models

type QType uint16

type DNSQuestion struct {
	QName  []string
	QType  QType
	QClass uint16
}

const (
	QTYPE_A     = 1  // A
	QTYPE_NS    = 2  // NS
	QTYPE_CNAME = 5  // CNAME
	QTYPE_SOA   = 6  // SOA
	QTYPE_TXT   = 16 // TXT
	QTYPE_MX    = 15 // MX
	QTYPE_AAAA  = 28 // AAAA
	QTYPE_PTR   = 12 // PTR
)

func (q QType) String() string {
	switch q {
	case QTYPE_A:
		return "A"
	case QTYPE_NS:
		return "NS"
	case QTYPE_CNAME:
		return "CNAME"
	case QTYPE_SOA:
		return "SOA"
	case QTYPE_TXT:
		return "TXT"
	case QTYPE_MX:
		return "MX"
	case QTYPE_AAAA:
		return "AAAA"
	case QTYPE_PTR:
		return "PTR"
	default:
		return "UNKNOWN"
	}
}
