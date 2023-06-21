package lastnbuffer_test

import (
	"testing"

	"github.com/eos175/lastnbuffer"
)

func TestLastNBuffer(t *testing.T) {
	bufferSize := 4
	buffer, err := lastnbuffer.NewLastNBuffer[int](bufferSize)
	if err != nil {
		t.Fatalf("Error creating buffer: %v", err)
	}

	// Test empty buffer
	elements := buffer.GetAll()
	if len(elements) != 0 {
		t.Errorf("Expected empty buffer, got: %v", elements)
	}

	// Test appending elements
	buffer.Append(1)
	buffer.Append(2)
	buffer.Append(3)
	buffer.Append(4)
	buffer.Append(5)
	buffer.Append(6)

	// Test GetLastN
	elements = buffer.GetLastN(4)
	expected := []int{3, 4, 5, 6}
	if !equalSlices(elements, expected) {
		t.Errorf("Expected last 3 elements to be %v, got: %v", expected, elements)
	}

	// Test ForEach
	var result []int
	buffer.ForEach(func(element int) bool {
		result = append(result, element)
		return true
	})
	if !equalSlices(result, expected) {
		t.Errorf("Expected ForEach result to be %v, got: %v", expected, result)
	}

	// Test Reset
	buffer.Reset()
	elements = buffer.GetAll()
	if len(elements) != 0 {
		t.Errorf("Expected empty buffer after Reset, got: %v", elements)
	}
}

// Helper function to compare two slices
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
