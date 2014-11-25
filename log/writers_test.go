package log

import (
	"bytes"
	"testing"
)

func TestNullWriter(t *testing.T) {
	n := new(NullWriter)

	if l, err := n.Write([]byte("1234567890")); l != 10 || err != nil {
		t.Log("Did expect to write 10 bytes and no error")
		t.Fail()
	}
}

func TestNewlineWriter(t *testing.T) {
	var buf []byte

	b := bytes.NewBuffer(buf)
	w := NewNewlineWriter(b)

	w.Write([]byte("123456789"))

	if b.Bytes()[9] != '\n' {
		t.Error("NewlineWriter must write a newline")
	}
}
