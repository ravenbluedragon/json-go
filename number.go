package jsongo

// Number represents a json number
// note: JSON is arbitrary precision but recommends float64
type Number float64

// json for type safety
func (Number) json() {}

type numberReader struct {
	char      rune
	err       error
	int_neg   bool
	int_part  uint64
	frac_part uint64
	frac_len  uint
	exp_neg   bool
	exp_part  uint64
}

var exponents [][]Number = [][]Number{
	{1, 1e1, 1e2, 1e4, 1e8, 1e16, 1e32, 1e64, 1e128, 1e256},
	{1, 1e-1, 1e-2, 1e-4, 1e-8, 1e-16, 1e-32, 1e-64, 1e-128, 1e-256},
}

func (nr *numberReader) number() (Number, error) {
	n := Number(nr.int_part)
	f := Number(nr.frac_part)
	if nr.int_neg {
		n = -n
		f = -f
	}
	if nr.frac_len > 17 {
		nr.err = PrecisionInsufficient
	}
	e := 0
	if nr.exp_neg {
		e -= int(nr.exp_part)
	} else {
		e += int(nr.exp_part)
	}
	if e <= -308 {
		nr.err = PrecisionInsufficient
	}
	exp := scale(e)
	return n*exp + f*scale(e-int(nr.frac_len)), nr.err
}

func scale(e int) Number {
	exps := exponents[0]
	if e < 0 {
		exps = exponents[1]
		e = -e
	}
	exp := Number(1)
	for i, j := 1, 1; i <= 256; i, j = i<<1, j+1 {
		if e&i != 0 {
			exp *= exps[j]
		}
	}
	return exp
}

// parseNumber reads a number from a reader
func parseNumber(r *reader) (Number, error) {
	nr := new(numberReader)
	for f := intSign; f != nil; {
		nr.char, nr.err = r.readRune()
		for a := nextFn; f != nil && a == nextFn; {
			f, a = f(nr)
			if a == invalidChar {
				return 0, InvalidCharacter{nr.char, r.position - 1}
			}
			if a == rewind {
				r.advance(-1)
			}
		}
	}
	if nr.err != nil {
		return 0, nr.err
	}
	return nr.number()
}

type action uint8

const (
	nextChar action = iota
	nextFn
	invalidChar
	lastAction
	rewind
)

type fn func(*numberReader) (fn, action)

func intSign(nr *numberReader) (fn, action) {
	if nr.err != nil {
		return nil, lastAction
	}
	if nr.char == '-' {
		nr.int_neg = true
		return intDigit, nextChar
	}
	return intDigit, nextFn
}

func intDigit(nr *numberReader) (fn, action) {
	if nr.err != nil {
		return nil, lastAction
	}
	switch nr.char {
	case '0':
		return intZero, nextChar
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return intDigits, nextFn
	case 'e', 'E':
		return expSign, nextChar
	default:
		return terminal, nextFn
	}
}

func intZero(nr *numberReader) (fn, action) {
	if nr.char == '.' && nr.err == nil {
		return fracDot, nextChar
	}
	return terminal, nextFn
}

func intDigits(nr *numberReader) (fn, action) {
	if nr.err != nil {
		return terminal, nextFn
	}
	switch nr.char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		nr.int_part = 10*nr.int_part + uint64(nr.char-'0')
		return intDigits, nextChar
	case '.':
		return fracDot, nextChar
	case 'e', 'E':
		return expSign, nextChar
	default:
		return nil, invalidChar
	}
}

func fracDot(nr *numberReader) (fn, action) {
	if nr.err != nil {
		return nil, lastAction
	}
	switch nr.char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return fracDigits, nextFn
	default:
		return nil, invalidChar
	}
}

func fracDigits(nr *numberReader) (fn, action) {
	if nr.err != nil {
		return terminal, nextFn
	}
	switch nr.char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if nr.frac_len < 18 {
			nr.frac_part = 10*nr.frac_part + uint64(nr.char-'0')
			nr.frac_len++
		}
		return fracDigits, nextChar
	case 'e', 'E':
		return expSign, nextChar
	default:
		return terminal, nextFn
	}
}

func expSign(nr *numberReader) (fn, action) {
	if nr.err != nil {
		return nil, lastAction
	}
	act := nextFn
	if nr.char == '-' {
		nr.exp_neg = true
		act = nextChar
	} else if nr.char == '+' {
		act = nextChar
	}
	return expDigit, act
}

func expDigit(nr *numberReader) (fn, action) {
	if nr.err != nil {
		return terminal, nextFn
	}
	switch nr.char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return expDigits, nextFn
	default:
		return nil, invalidChar
	}
}

func expDigits(nr *numberReader) (fn, action) {
	if nr.exp_part >= 308 && !nr.exp_neg {
		nr.err = NumberNotRepresentable
		return nil, lastAction
	}
	if nr.err != nil {
		return terminal, nextFn
	}
	switch nr.char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		nr.exp_part = 10*nr.exp_part + uint64(nr.char-'0')
		return expDigits, nextChar
	default:
		return terminal, nextFn
	}
}

func terminal(nr *numberReader) (fn, action) {
	if nr.err == UnexpectedEndOfDocument {
		nr.err = nil
		return nil, lastAction
	}
	switch nr.char {
	case ',', ']', '}', '\n', '\r', '\t', ' ':
		return nil, rewind
	default:
		return nil, invalidChar
	}
}
