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

// InvalidCharacter signals that there is no context this character is correct
type InvalidCharacter struct {
	character rune
	position  int
}

func (e InvalidCharacter) Error() string { return string(e.character) }
