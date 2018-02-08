package errors

import "fmt"

// Error represents a error interface all swagger framework errors implement
type Error interface {
    error
    Code() int
}

type MLCloudError struct {
    code    int
    message string
}

func (a *MLCloudError) Error() string {
    return a.message
}

func (a *MLCloudError) Code() int {
    return a.code
}

// New creates a new API error with a code and a message
func New(code int, message string, args ...interface{}) Error {
    if len(args) > 0 {
        return &MLCloudError{code, fmt.Sprintf(message, args...)}
    }
    return &MLCloudError{code, message}
}
