package msg

import (
	"testing"
)

func TestNewFilter(t *testing.T) {
	f := NewFilter()

	if len(f.in) != 0 {
		t.Error("In vector not initialized properly")
	}

	if len(f.out) != 0 {
		t.Error("Out vector not initialized properly")
	}

	f.In(1)

	if len(f.in) != 1 && f.in[0] != 1 {
		t.Error("In entry not set correctly")
	}

	f.Out(2)

	if len(f.out) != 1 && f.out[0] != 2 {
		t.Error("Out entry not set correctly")
	}
}

func TestFilter(t *testing.T) {
	m := Msg(1, []byte("Test message"))

	f := NewFilter()
	// Both filtered in and out, is not allowed.
	f.Out(1)
	f.In(1)

	if m.Allowed(f) {
		t.Error("Did not expect message to be allowed")
	}

	f = NewFilter()
	// Filtered in must be allowed.
	f.In(1)

	if !m.Allowed(f) {
		t.Error("Did not expect message to be disallowed")
	}

	f = NewFilter()
	// Filtered out, but not this one.
	f.Out(10)

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
