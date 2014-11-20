package log

import (
	"testing"
)

func TestNullWriter(t *testing.T) {
	n := new(NullWriter)

	if l, err := n.Write([]byte("1234567890")); l != 10 || err != nil {
		t.Log("Did expect to write 10 bytes and no error")
		t.Fail()
	}
}
