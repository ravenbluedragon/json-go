package jsongo

// A GeneralError is a simple string state that needs no additional context
type GeneralError string

func (e GeneralError) Error() string { return string(e) }

// UnexpectedEndOfDocument signals that a process tried to read beyond the final character in a document
const UnexpectedEndOfDocument = GeneralError("Unexpected End of Docuemnt")

// PeekBackwards signals an attempt to peek earlier in the document
const PeekBackwards = GeneralError("Cannot peek Backwards")
