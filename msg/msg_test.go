package msg

import (
	"testing"
	"time"
)

func TestString(t *testing.T) {
	m := Message{10, 1416926265, []byte("Some text here")}
	mstr := m.String()

	if mstr != "10 1416926265 Some text here" {
		t.Error("Unexpected string format for Message")
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

func TestTimestap(t *testing.T) {
	ts := makeTimestamp() / 1000
	tm := time.Unix(ts, 0)

	if tm.After(time.Now()) {
		t.Error("Time was given in the future")
	}

	if ts <= 1416926265 {
		t.Error("Time was given before the test was written")
	}
}

func TestCasting(t *testing.T) {
	str := "some random string"
	m := Cast([]byte(str))

	if string(m.Data()) != str {
		t.Error("Casting did not preserve data")
	}

	str = "some data"
	m = Msg(2, []byte(str))
	mc := Cast(m)

	if string(mc.Data()) != str || mc.senderId != 2 {
		t.Error("Casting did not preserve a Message")
	}
}

func TestInvalid(t *testing.T) {
	str := []byte("no string But string was expected")
	_, err := FromString(str)

	if err == nil {
		t.Error("Expected error after passing invalid message")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic()")
		}
	}()

	// Will trigger a panic.
	Cast(interface{}(NewFilter()))
}
