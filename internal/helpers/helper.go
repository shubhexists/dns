package helpers

// helper function to convert array of bytes to string
func ByteToString(data []byte) string {
	s := ""

	for _, b := range data {
		s = s + string(b)
	}
	return s
}

// array of function to convert array of bytes to int
func ByteToInt(data []byte) []int {
	i := []int{}
	for _, b := range data {
		i = append(i, int(b))
	}
	return i
}
