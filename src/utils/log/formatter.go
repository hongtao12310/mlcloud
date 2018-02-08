package log

// Formatter formats records in different ways: text, json, etc.
type Formatter interface {
	Format(*Record) ([]byte, error)
}
