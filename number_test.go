package jsongo

import "testing"

func TestNumber(t *testing.T) {
	type test struct {
		doc   string
		value Number
		pos   int
		err   error
	}
	table := []test{
		{"0", 0, 1, nil},
		{"00", 0, 2, InvalidCharacter{'0', 1}},
		{"000.0", 0, 2, InvalidCharacter{'0', 1}},
		{"0.0", 0, 3, nil},
		{"01.2", 0, 2, InvalidCharacter{'1', 1}},
		{"1234567890", 1234567890, 10, nil},
		{"1E400", 0, 5, NumberNotRepresentable},
		{"1e-400", 0, 6, PrecisionInsufficient},
		{"3.141592653589793238462643383279", 3.141592653589793, 32, PrecisionInsufficient},
		{"30.6e50", 30.6e50, 7, nil},
		{"30.6 e 50", 0, 6, InvalidCharacter{' ', 5}},
		{"-7.5", -7.5, 4, nil},
		{"-31e-5", -31e-5, 6, nil},
		{"-31e+5", -31e5, 6, nil},
		{"--7.5", 0, 2, InvalidCharacter{'-', 1}},
		{"-31e--5", 0, 6, InvalidCharacter{'-', 5}},
	}
	for _, tc := range table {
		r := &reader{tc.doc, 0}
		v, e := parseNumber(r)
		helper(t,
			log("Test case: %q", tc.doc),
			argument{"Value expected: %#v, received %#v", tc.value, v},
			argument{"Position expected: %d, received %d", tc.pos, r.position},
			argument{"Error expected: %v, received %v", tc.err, e},
		)
	}
}
