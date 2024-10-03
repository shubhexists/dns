package helpers

func ByteToString(data []byte) string {
	s := ""

	for _, b := range data {
		s = s + string(b)
	}
	return s
}

func ByteToInt(data []byte) []int {
	i := []int{}
	for _, b := range data {
		i = append(i, int(b))
	}
	return i
}
