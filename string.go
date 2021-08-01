package jsongo

// String represents a json string
type String string

// json for type safety
func (String) json() {}

// parseString reads a string from a reader
func parseString(*reader) (String, error) { return "", nil }
