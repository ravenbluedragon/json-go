package jsongo

import "testing"

func TestString(t *testing.T) {
	type test struct {
		doc   string
		value String
		pos   int
		err   error
	}
	table := []test{
		{"\"\"", "", 2, nil},
		{"   ", "", 3, OpeningQuoteMissing},
		{"   \"  ", "  ", 6, ClosingQuoteMissing},
		{"\"\"", "", 2, nil},
		{"\"\\", "", 2, UnexpectedEndOfDocument},
		{"\"a\"", "a", 3, nil},
		{"   \" \\t  a\"   ", " \t  a", 14, nil},
		{"\" \t  a\"   ", " ", 2, InvalidCharacter{'\t', 2}},
		{"\"\\ua2F5\"", "\ua2F5", 8, nil},
		{"\"\\u09AF\"", "\u09AF", 8, nil},
		{"\" \\\" \\\\ \\/ \\b \\f \\n \\r \\u09af \"", " \" \\ / \b \f \n \r \u09AF ", 31, nil},
		{"\" \\u04g  a\"   ", " ", 6, InvalidCharacter{'g', 6}},
		{"\" \\u04g", " ", 4, UnexpectedEndOfDocument},
	}
	for _, tc := range table {
		r := &reader{tc.doc, 0}
		v, e := parseString(r)
		helper(t,
			log("Test case: %q", tc.doc),
			argument{"Value expected: %q, received %q", tc.value, v},
			argument{"Position expected: %d, received %d", tc.pos, r.position},
			argument{"Error expected: %v, received %v", tc.err, e},
		)
	}
}
