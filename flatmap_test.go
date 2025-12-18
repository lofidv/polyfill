package polyfill

import "testing"

func TestFlatMap(t *testing.T) {
	words := []string{"hello world", "foo bar"}
	result := FlatMap(From(words), func(s string) []string {
		// Simple split simulation
		var parts []string
		current := ""
		for _, ch := range s {
			if ch == ' ' {
				if current != "" {
					parts = append(parts, current)
					current = ""
				}
			} else {
				current += string(ch)
			}
		}
		if current != "" {
			parts = append(parts, current)
		}
		return parts
	}).Slice()

	expected := []string{"hello", "world", "foo", "bar"}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %s at index %d, got %s", expected[i], i, result[i])
		}
	}
}

func TestFlatten(t *testing.T) {
	nested := [][]int{{1, 2}, {3, 4}, {5}}
	result := Flatten(From(nested)).Slice()

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, result[i])
		}
	}
}
