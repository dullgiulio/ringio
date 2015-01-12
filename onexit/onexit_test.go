package onexit

import (
	"testing"
)

func TestExit(t *testing.T) {
	exitVal := 0

	SetFunc(func(i int) {
		exitVal = i
	})

	Exit(10)

	if exitVal != 10 {
		t.Error("Did not exit with status 10")
	}
}

func TestDefer(t *testing.T) {
	exitVal := 0

	SetFunc(func(i int) {
		exitVal += i
	})

	Defer(func() { exitVal++ })
	Defer(func() { exitVal++ })

	Exit(1)

	if exitVal != 3 {
		t.Error("Some defer function was not called correctly")
	}
}

func TestPendingExit(t *testing.T) {
	exitVal := 0
	done := make(chan struct{})

	SetFunc(func(i int) {
		exitVal = i
		done <- struct{}{}
	})

	HandleInterrupt()
	PendingExit(10)

	<-done

	if exitVal != 10 {
		t.Error("Did not exit with status 10")
	}
}
