package store

// Driver defines methods that a configuration store driver must implement
type Driver interface {
	// Name returns a human-readable name of the driver
	Name() string
	// Read reads all the configurations from store
	Read() (map[string]interface{}, error)
	// Write writes the configurations to store, the configurations can be
	// part of all
	Write(map[string]interface{}) error
}
