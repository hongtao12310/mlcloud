package fs

import "testing"

func TestStatusCode(t *testing.T) {
    t.Logf("StatusDirectoryNotAFile: %d", StatusDirectoryNotAFile)
    if StatusFileNotFound != 1 {
        t.Fail()
    }

}

