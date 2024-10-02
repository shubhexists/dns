package helpers

import (
	"fmt"
	"testing"
)

func TestbyteToString(t *testing.T) {
	data := []byte{0x77, 0x77, 0x77}

	expected := "www"

	result := byteToString(data)

	if result != expected {
		fmt.Printf("Expected: %q Result: %q", expected, result)
	}
}
