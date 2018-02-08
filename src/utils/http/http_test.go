package http

import (
    "testing"
    "os"
)

func TestStat(t *testing.T) {
    _, err := os.Stat("test.txt")
    if err != nil && os.IsNotExist(err) {
        t.Error(err.Error())
    }
}

func TestCreateFile(t *testing.T)  {
    path := "test.txt"
    if _, err := os.Stat(path); os.IsNotExist(err) {
        _, err := os.Create(path)
        if err != nil {
            t.Errorf("failed create file: %s", path)
            return
        }

        os.Remove(path)
    }

}
