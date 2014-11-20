package log

import (
	"fmt"
	"io"
)

type NullWriter struct{}

func (n *NullWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

type NewlineWriter struct {
	w io.Writer
}

func NewNewlineWriter(w io.Writer) *NewlineWriter {
	return &NewlineWriter{w: w}
}

func (nw *NewlineWriter) Write(b []byte) (int, error) {
	return nw.w.Write([]byte(fmt.Sprintf("%s\n", b)))
}
