package utils

import (
    "strconv"
    "strings"
)

func ParseStringToInt(str string) (int, error) {
    if len(str) == 0 {
        return 0, nil
    }
    return strconv.Atoi(str)
}

func ParseStringToBool(str string) (bool, error) {
    return strings.ToLower(str) == "true" ||
        strings.ToLower(str) == "on", nil
}
