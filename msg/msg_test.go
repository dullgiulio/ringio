package msg

import (
	"testing"
)

func TestString(t *testing.T) {
	m := Msg(0, []byte("Test message text to parse"))
	if m.String() == "" {
		t.Error("Did not expect a message to be empty")
	}
}

func TestFromString(t *testing.T) {
	msg := []byte("12 1416213479589 Test message text to parse")
	m, err := FromString(msg)

	if err != nil {
		t.Error(err)
	}

	if m.senderId != 12 {
		t.Error("Expected senderId to be set correctly")
	}

	if m.time != 1416213479589 {
		t.Error("Expected time to be set correctly")
	}

	if string(m.data) != "Test message text to parse" {
		t.Error("Expected data to be set correctly")
	}
}
