package msg

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestNewFilter(t *testing.T) {
	f := NewFilter()

	if len(f.Included) != 0 {
		t.Error("In vector not initialized properly")
	}

	if len(f.Excluded) != 0 {
		t.Error("Out vector not initialized properly")
	}

	f.In(1)

	if len(f.Included) != 1 && f.Included[0] != 1 {
		t.Error("In entry not set correctly")
	}

	f.Out(2)

	if len(f.Excluded) != 1 && f.Excluded[0] != 2 {
		t.Error("Out entry not set correctly")
	}
}

func TestFilter(t *testing.T) {
	m := Msg(1, []byte("Test message"))

	if !m.Allowed(nil) {
		t.Error("Allowed returned false on nil filter")
	}

	f := NewFilter()
	// Both filtered in and out, is not allowed.
	f.Out(2)
	f.Out(1)
	f.In(2)
	f.In(1)

	if f.String() != "1,2,-1,-2" {
		t.Error("Unexpected string representation")
	}

	if m.Allowed(f) {
		t.Error("Did not expect message to be allowed")
	}

	f = NewFilter()
	// Filtered in must be allowed.
	f.In(1)

	if f.String() != "1" {
		t.Error("Unexpected string representation")
	}

	if !m.Allowed(f) {
		t.Error("Did not expect message to be disallowed")
	}

	f = NewFilter()
	// Filtered out, but not this one.
	f.Out(10)

	if f.String() != "-10" {
		t.Error("Unexpected string representation")
	}

	if !m.Allowed(f) {
		t.Error("Did not expect message to be disallowed")
	}

	f = NewFilter()
	// Something is filtered in, the rest must be filtered out.
	f.In(2)

	if m.Allowed(f) {
		t.Error("Did not expect message to be allowed")
	}
}

func TestCanGob(t *testing.T) {
	var buf bytes.Buffer

	f := NewFilter()
	f.In(1)
	f.Out(2)

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(f); err != nil {
		t.Error(err)
	}
}
