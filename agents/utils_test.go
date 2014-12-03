package agents

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/msg"
)

func TestWriteToChan(t *testing.T) {
	var buf bytes.Buffer

	for i := 0; i < 100; i++ {
		fmt.Fprintf(&buf, "Some string %d\n", i)
	}

	c := make(chan []byte)
	cancel := make(chan bool)

	go writeToChan(c, cancel, &buf)

	for i := 0; i < 100; i++ {
		// Won't contain the newline now.
		s := fmt.Sprintf("Some string %d", i)
		readS := <-c

		if s != string(readS) {
			t.Error("Unexpected difference in reading to channel")
			break
		}
	}
}

type BufferCloser bytes.Buffer

func (b BufferCloser) Close() error {
	return nil
}

func (b BufferCloser) Read(p []byte) (n int, err error) {
	buf := bytes.Buffer(b)
	return buf.Read(p)
}

func (b BufferCloser) Write(p []byte) (n int, err error) {
	buf := bytes.Buffer(b)
	return buf.Write(p)
}

func TestWriteToRingbuf(t *testing.T) {
	var buf bytes.Buffer

	for i := 0; i < 100; i++ {
		fmt.Fprintf(&buf, "Some string %d\n", i)
	}

	go readLogs(t)

	ring := ringbuf.NewRingbuf(100)
	go ring.Run()

	cancel := make(chan bool, 1)
	cancel <- true

	cancelled := writeToRingbuf(0, BufferCloser(buf), ring, cancel, nil)

	if !cancelled {
		t.Error("Cancelling writeToRingbuf didn't cancel")
	}

	ring.Eof()

	reader := ringbuf.NewReader(ring)
	readCh := reader.ReadCh()
	i := 0

	for data := range readCh {
		str := fmt.Sprintf("Some string %d", i)

		if string(data.([]byte)) != str {
			t.Error("Unexpected string was written")
		}

		i++
	}

	reader.Cancel()
	ring.Cancel()
}

func TestReadFromRingbuf(t *testing.T) {
	var buf bytes.Buffer

	ring := ringbuf.NewRingbuf(100)
	go ring.Run()

	for i := 0; i < 10; i++ {
		m := msg.Msg(0, []byte(fmt.Sprintf("Some string %d\n", i)))
		ring.Write(m)
	}

	ring.Eof()

	cancel := make(chan bool)

	cancelled := readFromRingbuf(BufferCloser(buf), ring, cancel, nil)

	if cancelled {
		t.Error("Did not expect readFromRingbuf to be cancelled")
	}

	ring.Cancel()
}
