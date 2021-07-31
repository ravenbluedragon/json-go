package jsongo

// A GeneralError is a simple string state that needs no additional context
type GeneralError string

func (e GeneralError) Error() string { return string(e) }

// UnexpectedEndOfDocument signals that a process tried to read beyond the final character in a document
const UnexpectedEndOfDocument = GeneralError("Unexpected End of Docuemnt")

// PeekBackwards signals an attempt to peek earlier in the document
const PeekBackwards = GeneralError("Cannot peek Backwards")

// KeywordNotFound signals that the reader does not match a keyword
const KeywordNotFound = GeneralError("No Keyword matched")

// num converts an int to a decimal string
func num(n int) string {
	if n == 0 {
		return "0"
	}
	left := 0
	var digs []rune
	if n < 0 {
		digs = append(digs, '-')
		n = -n
		left++
	}
	for n > 0 {
		d := rune(n % 10)
		digs = append(digs, '0'+d)
		n /= 10
	}
	right := len(digs) - 1
	for left < right {
		digs[left], digs[right] = digs[right], digs[left]
		left++
		right--
	}
	return string(digs)
}

// InvalidCharacter signals that there is no context this character is correct
type InvalidCharacter struct {
	character rune
	position  int
}

func (e InvalidCharacter) Error() string {
	return "Invalid Character " + string(e.character) + " at position " + num(e.position)
}
