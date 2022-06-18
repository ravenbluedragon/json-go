package jsongo

import "testing"

type logline struct {
	format    string
	arguments []interface{}
}

// log is a helper function to simplify testing
func log(format string, args ...interface{}) logline {
	return logline{format, args}
}

type argument struct {
	format   string
	expected interface{}
	received interface{}
}

// helper is a function to simplify table testing
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
		t.Log("-----")
	}
}
