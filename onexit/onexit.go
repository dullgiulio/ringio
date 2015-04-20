package onexit

import (
	"os"
	"os/signal"
	"sync"
)

type Func func()

type onExit struct {
	m  sync.Mutex
	i  int
	s  chan os.Signal
	fs []Func
	e  func(int)
}

var _onExit onExit

func init() {
	_onExit = onExit{
		s:  make(chan os.Signal),
		fs: make([]Func, 0),
		e:  os.Exit,
	}
}

func SetFunc(e func(int)) {
	_onExit.m.Lock()
	defer _onExit.m.Unlock()

	_onExit.e = e
}

func Defer(f Func) {
	_onExit.m.Lock()
	defer _onExit.m.Unlock()

	_onExit.fs = append(_onExit.fs, f)
}

func Exit(i int) {
	_onExit.m.Lock()
	defer _onExit.m.Unlock()

	for _, f := range _onExit.fs {
		f()
	}

	_onExit.e(i)
}

func PendingExit(i int) {
	_onExit.m.Lock()
	_onExit.i = i
	_onExit.m.Unlock()

	_onExit.s <- os.Interrupt
}

func HandleInterrupt() {
	signal.Notify(_onExit.s, os.Interrupt)

	go func() {
		<-_onExit.s

		_onExit.m.Lock()
		i := _onExit.i
		_onExit.m.Unlock()

		Exit(i)
	}()
}
