package helpers

import (
	"reflect"
	"testing"
)

func TestByteToString(t *testing.T) {
	data := []byte{0x77, 0x77, 0x77}

	expected := "www"

	result := ByteToString(data)

	if result != expected {
		t.Errorf("Expected: %q Result: %q", expected, result)
	}
}

func TestByteToInt(t *testing.T) {
	data := []byte{0x8E, 0xFA, 0x48, 0x64}

	expected := []int{142, 250, 72, 100}

	result := ByteToInt(data)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %+v and Got: %+v", expected, result)
	}
}
