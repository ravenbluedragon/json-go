package jsongo

import "testing"

func TestKeyword(t *testing.T) {
	type test struct {
		doc    string
		output Value
		pos    int
		err    error
	}
	_true := Bool(true)
	_false := Bool(false)
	table := []test{
		{"true", _true, 4, nil},
		{"false", _false, 5, nil},
		{"\ttrue\r\n", _true, 7, nil},
		{" null   ", nil, 8, nil},
		{"tr ue", nil, 0, KeywordNotFound},
		{"truth", nil, 0, KeywordNotFound},
		{"truest", _true, 4, InvalidCharacter{'s', 4}},
		{"true,", _true, 4, nil},
		{"true]   g", _true, 4, nil},
		{"true}", _true, 4, nil},
		{"true  }", _true, 6, nil},
	}
	for _, tc := range table {
		r := &reader{tc.doc, 0}
		v, e := parseKeyword(r)
		helper(t,
			log("Test case: %#v", tc.doc),
			argument{"Output expected %#v, Received %#v", tc.output, v},
			argument{"Position expected %d, Received %d", tc.pos, r.position},
			argument{"Error expected %v, received %v", tc.err, e},
		)
	}
}
