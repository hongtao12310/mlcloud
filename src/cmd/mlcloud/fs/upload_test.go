package fs

import (
    "testing"
    "strings"
)

func TestTrimLeft(t *testing.T) {
    s1 := "/Users/hongtaozhang/.mlcloud/config"
    s2 := "/Users/hongtaozhang/.mlcloud"

    t.Log(strings.TrimPrefix(s1, s2))

}
