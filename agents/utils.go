package agents

import (
	"bufio"
	"io"
	"sync"

	"bitbucket.org/dullgiulio/ringbuf"
	"bitbucket.org/dullgiulio/ringio/log"
	"bitbucket.org/dullgiulio/ringio/msg"
)

func _writeToChan(c chan<- []byte, cancel chan bool, reader io.ReadCloser) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		c <- scanner.Bytes()
	}

	if err := scanner.Err(); err != nil {
		log.Error(log.FacilityAgent, err)
	}

	log.Debug(log.FacilityAgent, "Writing into channel from input has terminated")

	reader.Close()

	close(c)
}

func writeToRingbuf(id int, reader io.ReadCloser, ring *ringbuf.Ringbuf, cancel chan bool, wg *sync.WaitGroup) (cancelled bool) {
	if wg != nil {
		defer wg.Done()
	}

	c := make(chan []byte)

	go _writeToChan(c, cancel, reader)

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
				log.Error(log.FacilityAgent, err)
				return
			}
			if err := output.WriteByte('\n'); err != nil {
				log.Error(log.FacilityAgent, err)
				return
			}
			if err := output.Flush(); err != nil {
				log.Error(log.FacilityAgent, err)
				return
			}
		case <-cancel:
			cancelled = true
			log.Debug(log.FacilityAgent, "Read from ringbuf has been cancelled")
			return
		}
	}
}

func readFromRingbuf(writer io.WriteCloser, ring *ringbuf.Ringbuf, cancel <-chan bool, wg *sync.WaitGroup) (cancelled bool) {
	if wg != nil {
		defer wg.Done()
	}

	reader := ringbuf.NewRingbufReader(ring)
	output := bufio.NewWriter(writer)
	c := reader.ReadCh()

	cancelled = _readInnerLoop(c, cancel, output)

	reader.Cancel()
	writer.Close()

	return
}
