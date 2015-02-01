package msg

import (
	"fmt"
	"time"
)

type Flags int

const (
	MSG_FORMAT_NULL Flags = 1 << iota
	MSG_FORMAT_ID
	MSG_FORMAT_TIMESTAMP
	MSG_FORMAT_DATA
	MSG_FORMAT_NEWLINE
	MSG_FORMAT_META = MSG_FORMAT_ID | MSG_FORMAT_TIMESTAMP | MSG_FORMAT_DATA
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

// XXX: We implement only the combinations we actually use.
func (m Message) Format(f Flags) string {
	mask := MSG_FORMAT_META | MSG_FORMAT_NEWLINE

	if (f & mask) == mask {
		return fmt.Sprintf("%d %d %s\n", m.senderID, m.time, m.data)
	}

	mask = MSG_FORMAT_META

	if (f & mask) == mask {
		return fmt.Sprintf("%d %d %s", m.senderID, m.time, m.data)
	}

	if (f & MSG_FORMAT_NEWLINE) == MSG_FORMAT_NEWLINE {
		return fmt.Sprintf("%s\n", m.data)
	}

	return fmt.Sprintf("%s", m.data)
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
