// Package jsongo is a learning exercise on the json format
package jsongo

// Bool represents a json boolean true / false
type Bool bool

// json for type safety
func (Bool) json() {}

// parseKeyword reads a keyword from a reader
func parseKeyword(*reader) (Value, error) { return nil, nil }
