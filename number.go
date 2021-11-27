package jsongo

// Number represents a json number
// note: JSON is arbitrary precision but recommends float64
type Number float64

// json for type safety
func (Number) json() {}

// parseNumber reads a number from a reader
func parseNumber(*reader) (Number, error) {
	return 0, nil
}
