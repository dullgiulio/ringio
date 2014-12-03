package agents

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/log"
	"github.com/dullgiulio/ringio/msg"
)

// TODO: Must make this terminate somehow.
func writeToChan(c chan<- []byte, cancel chan bool, reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		c <- scanner.Bytes()
	}

	if err := scanner.Err(); err != nil {
		log.Error(log.FacilityAgent, fmt.Errorf("bufio.Scanner: %v", err))
	}

	log.Debug(log.FacilityAgent, "Writing into channel from input has terminated")

	close(c)
}

func writeToRingbuf(id int, reader io.ReadCloser, ring *ringbuf.Ringbuf, cancel chan bool, wg *sync.WaitGroup) (cancelled bool) {
	if wg != nil {
		defer wg.Done()
	}

	c := make(chan []byte)

	go func() {
		writeToChan(c, cancel, reader)
		reader.Close()
	}()

	for {
		select {
		case data, ok := <-c:
			if !ok {
				return
			}

			ring.Write(msg.Msg(id, data))
		case <-cancel:
			cancelled = true
			log.Debug(log.FacilityAgent, "Writing into ringbuf from input has been cancelled")
			reader.Close()
			return
		}
	}
}

func _readInnerLoop(c <-chan interface{}, cancel <-chan bool, output *bufio.Writer) (cancelled bool) {
	for {
		select {
		case data := <-c:
			if data == nil {
				return
			}

			m := msg.Cast(data)

			if _, err := output.Write(m.Data()); err != nil {
				log.Error(log.FacilityAgent, fmt.Errorf("bufio.Write: %v", err))
				return
			}
			if err := output.WriteByte('\n'); err != nil {
				log.Error(log.FacilityAgent, fmt.Errorf("bufio.WriteByte: %v", err))
				return
			}
			if err := output.Flush(); err != nil {
				log.Error(log.FacilityAgent, fmt.Errorf("bufio.Flush: %v", err))
				return
			}
		case <-cancel:
			cancelled = true
			return
		}
	}
}

func readFromRingbuf(writer io.WriteCloser, ring *ringbuf.Ringbuf, cancel <-chan bool, wg *sync.WaitGroup) (cancelled bool) {
	if wg != nil {
		defer wg.Done()
	}

	reader := ringbuf.NewReader(ring)
	output := bufio.NewWriter(writer)
	c := reader.ReadCh()

	cancelled = _readInnerLoop(c, cancel, output)

	reader.Cancel()
	writer.Close()

	log.Debug(log.FacilityAgent, "Read from ringbuf has been cancelled")

	return
}
