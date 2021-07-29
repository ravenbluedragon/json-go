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
