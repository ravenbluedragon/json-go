package jsongo

import "testing"

type logline struct {
	format    string
	arguments []interface{}
}

func log(format string, args ...interface{}) logline {
	return logline{format, args}
}

type argument struct {
	format   string
	expected interface{}
	received interface{}
}

func helper(t *testing.T, log logline, args ...argument) {
	fail := false
	for _, a := range args {
		if a.expected != a.received {
			t.Errorf(a.format, a.expected, a.received)
			fail = true
		}
	}
	if fail {
		t.Logf(log.format, log.arguments...)
	}
}

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
		helper(t,
			log("Test case: doc %#v, start %d", tc.doc, tc.start),
			argument{"Expected '%c', Received '%c'", tc.out.value, c},
			argument{"Expected %v, Received %v", tc.out.err, e},
			argument{"Expected %d, Received %d", tc.out.pos, r.position},
		)
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
		helper(t,
			log("Peek %d bytes", tc.count),
			argument{"Reader has been modified: %d, %d", pos, r.position},
			argument{"Expected %#v, Received %#v", tc.value, v},
			argument{"Expected %v, Received %v", tc.err, e},
		)
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
		{"true", 0, 4, 4, nil},
		{"false", 0, 5, 5, nil},
		{"asdf", 1, -2, 1, UnexpectedEndOfDocument},
		{"asdf", 1, -1, 0, nil},
		{"asdf", 1, 0, 1, nil},
		{"asdf", 1, 1, 2, nil},
		{"asdf", 1, 2, 3, nil},
		{"asdf", 1, 3, 4, nil},
		{"asdf", 1, 4, 1, UnexpectedEndOfDocument},
	}
	for _, tc := range table {
		r := reader{tc.doc, tc.pos}
		e := r.advance(tc.count)
		helper(t,
			log("Advance %d bytes, doc %#v, pos %d", tc.count, tc.doc, tc.pos),
			argument{"Expected: %d, Received %d", tc.target, r.position},
			argument{"Expected %v, Received %v", tc.err, e},
		)
	}
}

func TestReaderSkipWhitespace(t *testing.T) {
	table := []struct {
		doc    string
		pos    int
		target int
	}{
		{"", 0, 0},
		{"", 7, 7},
		{"a", 0, 0},
		{"a b ", 0, 0},
		{"a b ", 1, 2},
		{"a b ", 2, 2},
		{"a b ", 3, 4},
		{"a b ", 4, 4},
		{"a\t\tc\r", 1, 3},
		{"a\t\tc\r", 3, 3},
		{"a\t\tc\r", 4, 5},
		{"\t\r\n     ", 2, 8},
	}
	for _, tc := range table {
		r := reader{tc.doc, tc.pos}
		r.skipWhitespace()
		helper(t,
			log("Skip whitespace in doc %#v from %d", tc.doc, tc.pos),
			argument{"Expected %d, Received %d", tc.target, r.position},
		)
	}
}
