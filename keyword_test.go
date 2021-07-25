package jsongo

import "testing"

func TestNull(t *testing.T) {
	if val := Parse("null"); val != nil {
		t.Fatalf("Expected: %v; Received %v", nil, val)
	}
}
