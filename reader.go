package jsongo

// the reader maintains the position on the document and has several useful methods
type reader struct {
	document string
	position int
}

// readRune will read the next Rune from the document and advance the position
func (*reader) readRune() (rune, error) { return 0, nil }
