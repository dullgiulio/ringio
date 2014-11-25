package config

import (
	"testing"
)

func TestDefaults(t *testing.T) {
	if C.RingbufSize != defaults.RingbufSize {
		t.Error("Invalid default")
	}

	if C.RingbufLogSize != defaults.RingbufLogSize {
		t.Error("Invalid default")
	}

	if C.MaxLineSize != defaults.MaxLineSize {
		t.Error("Invalid default")
	}

	if C.AutoRun != defaults.AutoRun {
		t.Error("Invalid default")
	}

	if C.AutoExit != defaults.AutoExit {
		t.Error("Invalid default")
	}

	if C.AutoLock != defaults.AutoLock {
		t.Error("Invalid default")
	}

	if C.PrintLog != defaults.PrintLog {
		t.Error("Invalid default")
	}
}

func TestRingbufInit(t *testing.T) {
	if GetLogRingbuf() != nil {
		t.Error("Ringbuf allocated before being initialized")
	}

	Init()

	GetLogRingbuf().Cancel()
}
