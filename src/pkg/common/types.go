package common

// define the system config type
type SYSConfig interface {
    Get(key string) interface{}
    Load() error
}