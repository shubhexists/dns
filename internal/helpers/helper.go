package helpers

func byteToString(data []byte) string {
	s := ""

	for _, b := range data {
		s = s + string(b)
	}
	return s
}
