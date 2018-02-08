// provide assert utils for the unit test assert statement
package utils

import (
    "fmt"
    "testing"
)

func Equal(t *testing.T, a interface{}, b interface{}) {
    if a == b {
        return
    }
    message := fmt.Sprintf("%v != %v", a, b)
    t.Fatal(message)
}

// if err is not nil fatal
func Nil(t *testing.T, err interface{}) {
    if err == nil {
        return
    }

    message := "err is not nil"
    t.Fatal(message)
}