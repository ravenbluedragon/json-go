package jsongo

// the reader maintains the position on the document and has several useful methods
type reader struct {
	document string
	position int
}

// readRune will read the next Rune from the document and advance the position
func (r *reader) readRune() (rune, error) {
	if r.position >= len(r.document) {
		return 0, UnexpectedEndOfDocument
	}
	c := r.document[r.position]
	r.position++
	return rune(c), nil
}

// peek will return the next count bytes of the string without changing position
func (r reader) peek(count int) (string, error) {
	if count < 0 {
		return "", PeekBackwards
	}
	if r.position+count > len(r.document) {
		return "", UnexpectedEndOfDocument
	}
	return r.document[r.position : r.position+count], nil
}

// advance will move the position forward by count
func (r *reader) advance(count int) error {
	if count == 0 {
		return nil
	}
	t := r.position + count
	if t < 0 || t > len(r.document) {
		return UnexpectedEndOfDocument
	}
	r.position = t
	return nil
}

// skipWhitespace will advance past all whitespace
func (r *reader) skipWhitespace() {
	for {
		c, e := r.readRune()
		if e == UnexpectedEndOfDocument {
			return
		}
		switch c {
		case ' ', '\t', '\r', '\n':
			continue
		}
		_ = r.advance(-1) // NOTE: an UnexpectedEndOfDocument error would be ok here
		return
	}
}
