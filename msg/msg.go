package msg

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Format int

const (
	FORMAT_NULL Format = 1 << iota
	FORMAT_ID
	FORMAT_TIMESTAMP
	FORMAT_DATA
	FORMAT_NEWLINE
	FORMAT_META = FORMAT_ID | FORMAT_TIMESTAMP | FORMAT_DATA
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
func (m Message) WriteFormat(w io.Writer, f Format) (int, error) {
	mask := FORMAT_META | FORMAT_NEWLINE

	if (f & mask) == mask {
		return fmt.Fprintf(w, "%s: % 3d: %s\n", time.Unix(m.time, 0).Format(time.RFC3339), m.senderID, m.data)
	}

	mask = FORMAT_META

	if (f & mask) == mask {
		return fmt.Fprintf(w, "%s: % 3d: %s", time.Unix(m.time, 0).Format(time.RFC3339), m.senderID, m.data)
	}

	if (f & FORMAT_NEWLINE) == FORMAT_NEWLINE {
		return fmt.Fprintf(w, "%s\n", m.data)
	}

	return fmt.Fprintf(w, "%s", m.data)
}

func (m Message) Format(f Format) string {
	var b bytes.Buffer

	m.WriteFormat(&b, f)
	return b.String()
}

func (m Message) String() string {
	return fmt.Sprintf("%s: % 3d: %s", time.Unix(m.time, 0).Format(time.RFC3339), m.senderID, m.data)
}

func (m Message) Data() []byte {
	return m.data
}

func Parse(msg []byte) (Message, error) {
	m := Message{}

	if len(msg) < 34 {
		return m, fmt.Errorf("Invalid string to parse: string too short")
	}

	if t, err := time.Parse(time.RFC3339, string(msg[0:25])); err != nil {
		return m, err
	} else {
		m.time = t.Unix()
	}

	str := string(msg[27:])
	end := strings.Index(str, ": ")

	if end < 0 {
		return m, fmt.Errorf("Invalid string to parse: agent ID not found")
	}

	if senderid, err := strconv.Atoi(strings.TrimSpace(str[0:end])); err != nil {
		return m, err
	} else {
		m.senderID = senderid
	}

	m.data = []byte(str[end+2:])

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
