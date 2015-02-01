package msg

import (
	"testing"
	"time"
)

func TestString(t *testing.T) {
	m := Message{10, 1416926265, []byte("Some text here")}
	mstr := m.String()

	if mstr != "2014-11-25T16:37:45+02:00:  10: Some text here" {
		t.Error("Unexpected string format for Message")
	}
}

func TestFormat(t *testing.T) {
	m := Message{10, 1416926265, []byte("Some text here")}

	if m.Format(FORMAT_META|FORMAT_NEWLINE) != "2014-11-25T16:37:45+02:00:  10: Some text here\n" {
		t.Error("Unexpected formatted message with newline")
	}

	if m.Format(FORMAT_META) != "2014-11-25T16:37:45+02:00:  10: Some text here" {
		t.Error("Unexpected fromatted message")
	}

	if m.Format(FORMAT_NEWLINE) != "Some text here\n" {
		t.Error("Unexpected data message with newline")
	}

	if m.Format(FORMAT_DATA) != "Some text here" {
		t.Error("Unexpected data message")
	}
}

func TestParse(t *testing.T) {
	msg := []byte("2014-11-25T16:37:45+02:00:   5: Test message text to parse")
	m, err := Parse(msg)

	if err != nil {
		t.Error(err)
	}

	if m.senderID != 5 {
		t.Error("Expected senderID to be set correctly")
	}

	if m.time != 1416926265 {
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

	if string(mc.Data()) != str || mc.senderID != 2 {
		t.Error("Casting did not preserve a Message")
	}
}

func TestInvalid(t *testing.T) {
	str := []byte("no string but string was expected")
	_, err := Parse(str)

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
