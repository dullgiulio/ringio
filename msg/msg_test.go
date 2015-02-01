package msg

import (
	"fmt"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	m := Message{10, 1416926265, []byte("Some text here")}
	tstr := time.Unix(1416926265, 0).Format(time.RFC3339)
	mstr := m.String()

	if mstr != fmt.Sprintf("%s:  10: Some text here", tstr) {
		t.Error(fmt.Errorf("Unexpected string format for Message, got '%s'", mstr))
	}
}

func TestFormat(t *testing.T) {
	m := Message{10, 1416926265, []byte("Some text here")}
	tstr := time.Unix(1416926265, 0).Format(time.RFC3339)
	mstr := m.Format(FORMAT_META | FORMAT_NEWLINE)

	if mstr != fmt.Sprintf("%s:  10: Some text here\n", tstr) {
		t.Error(fmt.Errorf("Unexpected formatted message with newline, got '%s'", mstr))
	}

	mstr = m.Format(FORMAT_META)
	if mstr != fmt.Sprintf("%s:  10: Some text here", tstr) {
		t.Error(fmt.Errorf("Unexpected formatted message, got '%s'", mstr))
	}

	mstr = m.Format(FORMAT_NEWLINE)
	if mstr != "Some text here\n" {
		t.Error(fmt.Errorf("Unexpected data message with newline, got '%s'", mstr))
	}

	mstr = m.Format(FORMAT_DATA)
	if mstr != "Some text here" {
		t.Error(fmt.Errorf("Unexpected data message, got '%s'", mstr))
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
