package pipe

import (
	"os"
	"testing"

	"github.com/dullgiulio/ringio/log"
)

func TestCreatePipe(t *testing.T) {
	pipeName := "/tmp/testp"

	go func() {
		nw := new(log.NullWriter)
		log.AddWriter(nw)

		if !log.Run(log.LevelError) {
			t.Log("Did not expect logger to return an error")
			t.Fail()
		}
	}()

	p := New(pipeName)
	if err := p.Create(); err != nil {
		t.Error(err)
	}
	if ok := p.OpenWrite(); !ok {
		t.Log("Did not expect OpenWrite to fail")
		t.Fail()
	} else {
		defer os.Remove(pipeName)

		if stat, err := p.file.Stat(); err != nil {
			t.Error(err)
			t.Fail()
		} else {
			mode := stat.Mode() & os.ModeNamedPipe
			if mode != mode|os.ModeNamedPipe {
				t.Log("Expected file to be a named pipe")
				t.Fail()
			}
		}

		if i, err := p.Write([]byte("something")); i < 5 || err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	pr := New(pipeName)
	if ok := pr.OpenRead(); !ok {
		t.Log("Did not expect OpenRead to fail")
		t.Fail()
	} else {
		os.Remove(pipeName)

		var buf [100]byte

		if n, err := p.Read(buf[:]); n < 5 || err != nil {
			t.Log("Read: ", n, "Error: ", err)
			t.Fail()
		}
	}

	p.Close()
}
