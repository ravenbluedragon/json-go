// Package jsongo is a learning exercise on the json format
package jsongo

// Bool represents a json boolean true / false
type Bool bool

// json for type safety
func (Bool) json() {}

var keywords = map[string]Value{"true": Bool(true), "false": Bool(false), "null": nil}

// parseKeyword reads a keyword from a reader
func parseKeyword(r *reader) (Value, error) {
	r.skipWhitespace()
	for k, v := range keywords {
		l := len(k)
		p, e := r.peek(l)
		if e != nil || p != k {
			continue
		}
		_ = r.advance(l) // if peek succeeded advance must too
		return v, checkValidPost(r)
	}
	return nil, KeywordNotFound
}

// checkValidPost checks whether that character after a keyword is valid and skips whitespace if needed
func checkValidPost(r *reader) error {
	p, e := r.peek(1)
	if e != nil {
		// if there is nothing to peek, there is no character to worry about
		return nil
	}
	switch c := rune(p[0]); c {
	case ' ', '\t', '\r', '\n':
		r.skipWhitespace()
	case ',', ']', '}':
	default:
		return InvalidCharacter{c, r.position}
	}
	return nil
}
