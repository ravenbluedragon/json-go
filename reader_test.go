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

func TestReaderPeek(t *testing.T) {
	pos := 1
	r := reader{"sample doc", pos}
	table := []struct {
		count int
		value string
		err   error
	}{
		{0, "", nil},
		{1, "a", nil},
		{6, "ample ", nil},
		{9, "ample doc", nil},
		{10, "", UnexpectedEndOfDocument},
		{15, "", UnexpectedEndOfDocument},
		{-1, "", PeekBackwards},
	}
	for _, tc := range table {
		v, e := r.peek(tc.count)
		t.Logf("Peek %d bytes", tc.count)
		if r.position != pos {
			t.Fatalf("Reader has been modified: %#v, %d", r.document, r.position)
		}
		if tc.value != v {
			t.Fatalf("Expected %#v, Received %#v", tc.value, v)
		}
		if tc.err != e {
			t.Fatalf("Expected %v, Received %v", tc.err, e)
		}
	}
}

func TestReaderAdvance(t *testing.T) {
	table := []struct {
		doc    string
		pos    int
		count  int
		target int
		err    error
	}{
		{"", 0, 0, 0, nil},
		{"", 0, 1, 0, UnexpectedEndOfDocument},
		{"", 0, -1, 0, UnexpectedEndOfDocument},
		{"asdf", 1, -2, 1, UnexpectedEndOfDocument},
		{"asdf", 1, -1, 0, nil},
		{"asdf", 1, 0, 1, nil},
		{"asdf", 1, 1, 2, nil},
		{"asdf", 1, 2, 3, nil},
		{"asdf", 1, 3, 1, UnexpectedEndOfDocument},
	}
	for _, tc := range table {
		r := reader{tc.doc, tc.pos}
		e := r.advance(tc.count)
		t.Logf("Advance %d bytes, %v", tc.count, r)
		if r.position != tc.target {
			t.Fatalf("Expected: %d, Received %d", tc.target, r.position)
		}
		if tc.err != e {
			t.Fatalf("Expected %v, Received %v", tc.err, e)
		}
	}
}
