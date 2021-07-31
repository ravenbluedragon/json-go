package jsongo

import "testing"

func TestGeneralError(t *testing.T) {
	table := []string{
		"a",
		"asldfk",
		"  asd e",
	}
	for _, text := range table {
		if s := GeneralError(text).Error(); s != text {
			t.Errorf("General Error expected %s, received %s", text, s)
		}
	}
}

func TestInvalidCharacter(t *testing.T) {
	table := []struct {
		char rune
		pos  int
		text string
	}{
		{'s', 5, "Invalid Character s at position 5"},
		{'a', 0, "Invalid Character a at position 0"},
		{'r', -10, "Invalid Character r at position -10"},
		{'6', 5000, "Invalid Character 6 at position 5000"},
	}
	for _, tc := range table {
		e := InvalidCharacter{tc.char, tc.pos}
		helper(t,
			log("Case '%c' at %d", tc.char, tc.pos),
			argument{"Expected: %s, Received %s", tc.text, e.Error()},
		)
	}
}
