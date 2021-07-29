package jsongo

import "testing"

func TestReaderRune(t *testing.T) {
	type output struct {
		value rune
		err   error
		pos   int
	}
	table := []struct {
		doc   string
		start int
		out   output
	}{
		{"", 0, output{0, UnexpectedEndOfDocument, 0}},
		{" ", 0, output{' ', nil, 1}},
		{"sadaf;lkj\\b", 5, output{';', nil, 6}},
		{"sadaf;lkj\\b", 10, output{'b', nil, 11}},
		{"sadaf;lkj\\b", 15, output{0, UnexpectedEndOfDocument, 15}},
		{" \r\t\n{}", 2, output{'\t', nil, 3}},
	}
	for _, tc := range table {
		r := reader{tc.doc, tc.start}
		c, e := r.readRune()
		if tc.out != (output{c, e, r.position}) {
			t.Logf("Test case: doc %#v, start %d", tc.doc, tc.start)
			if tc.out.value != c {
				t.Errorf("Expected '%c', Received '%c'", tc.out.value, c)
			}
			if tc.out.err != e {
				t.Errorf("Expected %v, Received %v", tc.out.err, e)
			}
			if tc.out.pos != r.position {
				t.Errorf("Expected %d, Received %d", tc.out.pos, r.position)
			}
		}
	}
}
