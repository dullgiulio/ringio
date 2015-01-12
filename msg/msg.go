package msg

import (
	"fmt"
	"time"
)

type Message struct {
	senderID int
	time     int64
	data     []byte
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Msg(senderID int, data []byte) Message {
	return Message{senderID, makeTimestamp(), data}
}

func (m Message) String() string {
	return fmt.Sprintf("%d %d %s", m.senderID, m.time, m.data)
}

func (m Message) Data() []byte {
	return m.data
}

func FromString(msg []byte) (Message, error) {
	m := Message{}
	ws := 0
	l := len(msg)
	metastr := ""

	for i := 0; i < l; i++ {
		if msg[i] == ' ' {
			ws++
		}

		if ws > 1 {
			m.data = msg[i+1:]
			metastr = string(msg[:i])
			break
		}
	}

	if metastr != "" {
		if _, err := fmt.Sscanf(metastr, "%d %d", &m.senderID, &m.time); err != nil {
			return m, err
		}
	}

	return m, nil
}

func Cast(i interface{}) Message {
	if m, ok := i.(Message); !ok {
		if d, ok := i.([]byte); !ok {
			panic("Cast to msg.Message failed")
		} else {
			return Msg(0, d)
		}
	} else {
		return m
	}
}
