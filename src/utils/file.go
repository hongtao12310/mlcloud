package utils

import (
    "os"
)

func MkdirIfNotExist(dir string, mode os.FileMode)  {
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        os.Mkdir(dir, mode)
    }
}