package jsongo

// String represents a json string
type String string

// json for type safety
func (String) json() {}

// parseString reads a string from a reader
func parseString(r *reader) (String, error) {
	z := String("")
	var runes []rune

	// find opening quote
	r.skipWhitespace()
	c, e := r.readRune()
	if e == UnexpectedEndOfDocument {
		return z, OpeningQuoteMissing
	}
	if e != nil {
		return z, e
	}

	done := false
	for !done {
		c, done, e = parseStringRune(r)
		if e != nil {
			return String(runes), e
		}
		if !done {
			runes = append(runes, c)
		}
	}
	r.skipWhitespace()
	return String(runes), nil
}

// parseStringRune is a helper to plan each rune of the string
func parseStringRune(r *reader) (rune, bool, error) {
	c, e := r.readRune()
	if e == UnexpectedEndOfDocument {
		return c, false, ClosingQuoteMissing
	}
	if c < 20 {
		return c, false, InvalidCharacter{c, r.position}
	}
	if c == '"' { // found closing quote
		return c, true, nil
	}
	if c == '\\' {
		c, e = readEscaped(r)
	}
	return c, false, e
}

var escapes = map[rune]rune{'t': '\t', 'b': '\b', 'f': '\f', 'r': '\r', 'n': '\n'}

// readEscaped reads the next escaped character sequence
func readEscaped(r *reader) (rune, error) {
	c, e := r.readRune()
	if e != nil {
		return 0, e
	}
	if c == 'u' {
		return readHex(r)
	}
	switch c {
	case '"', '\\', '/':
		return c, nil
	}
	if c, ok := escapes[c]; ok {
		return c, nil
	}
	return c, InvalidCharacter{c, r.position}
}

// readHex converts 4 hex digits to a unicode rune
func readHex(r *reader) (rune, error) {
	var v rune
	p, e := r.peek(4)
	if e != nil {
		return 0, e
	}
	for i, d := range p {
		v <<= 4
		switch d {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			v += (d - '0')
		case 'a', 'b', 'c', 'd', 'e', 'f':
			v += (d - 'a' + 10)
		case 'A', 'B', 'C', 'D', 'E', 'F':
			v += (d - 'A' + 10)
		default:
			_ = r.advance(i) // peek guarantees advance
			return 0, InvalidCharacter{d, r.position}
		}
	}
	_ = r.advance(4) // peek guarantees advance
	return v, nil
}
