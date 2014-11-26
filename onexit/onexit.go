package onexit

import (
	"os"
	"os/signal"
)

type OnExitFunc func()

type onExit struct {
	i  int
	s  chan os.Signal
	fs []OnExitFunc
	e  func(int)
}

var _onExit onExit

func init() {
	_onExit = onExit{
		s:  make(chan os.Signal),
		fs: make([]OnExitFunc, 0),
		e:  os.Exit,
	}
}

func SetExitFunc(e func(int)) {
	_onExit.e = e
}

func Defer(f OnExitFunc) {
	_onExit.fs = append(_onExit.fs, f)
}

func Exit(i int) {
	for _, f := range _onExit.fs {
		f()
	}

	_onExit.e(i)
}

func PendingExit(i int) {
	// This is not thread safe on purpose.
	_onExit.i = i
	_onExit.s <- os.Interrupt
}

func HandleInterrupt() {
	signal.Notify(_onExit.s, os.Interrupt)

	go func() {
		<-_onExit.s
		Exit(_onExit.i)
	}()
}
