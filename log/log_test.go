package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestLevel(t *testing.T) {
	if LevelDebug.String() != "DEBUG" {
		t.Error("Unexpected string value")
	}

	if l, e := LevelFromString("debug"); e != nil || l != LevelDebug {
		t.Error("Unexpected value from string")
	}

	if _, e := LevelFromString("no-level"); e == nil {
		t.Error("Expected error on invalid level string value")
	}
}

func TestLogging(t *testing.T) {
	var buf []byte

	b := bytes.NewBuffer(buf)
	w := NewNewlineWriter(b)

	AddWriter(w)

	proceed := make(chan struct{})

	go func() {
		if !Run(LevelWarn) {
			t.Error("Expected success after cancel")
		}
		proceed <- struct{}{}
	}()

	Debug(FacilityDefault, "Some debug message")
	Info(FacilityDefault, "Some info message")
	Cancel()

	<-proceed

	if b.String() != "" {
		t.Error("Debug or info message written even when minimum is Warn.")
	}

	b.Reset()

	go func() {
		if Run(LevelWarn) {
			t.Error("Expected run to return false after calling Fatal")
		}
		proceed <- struct{}{}
	}()

	Warn(FacilityDefault, "Some warn message")
	Fatal(FacilityDefault, "Some fatal error")

	<-proceed

	s := b.String()

	if !strings.Contains(s, "WARN: ringio: Some warn message") ||
		!strings.Contains(s, "FAIL: ringio: Some fatal error") {
		t.Error("Not all messages were correctly logged")
	}
}
