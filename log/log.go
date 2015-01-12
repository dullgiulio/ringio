package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	loggers []io.Writer
	lock    *sync.Mutex
	c       chan *Message
}

var _logger Logger

type Facility string

const (
	FacilityRing    Facility = "ringbuf"
	FacilityAgent            = "agent"
	FacilityDefault          = "ringio"
	FacilityStdout           = "stdout"
	FacilityStderr           = "stderr"
	FacilityPipe             = "pipe"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFail
	LevelCancel
)

type Message struct {
	facility Facility
	level    Level
	message  string
}

var prefixes = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFail:  "FAIL",
}

func (l Level) String() string {
	return prefixes[l]
}

func LevelFromString(s string) (Level, error) {
	s = strings.ToUpper(s)

	for l, sl := range prefixes {
		if sl == s {
			return l, nil
		}
	}

	return LevelInfo, errors.New("Invalid value")
}

func init() {
	_logger = Logger{
		loggers: make([]io.Writer, 0),
		lock:    new(sync.Mutex),
		c:       make(chan *Message),
	}
}

func _string(ds []interface{}) string {
	buffer := new(bytes.Buffer)

	for _, d := range ds {
		buffer.WriteString(fmt.Sprintf("%v ", d))
	}

	return buffer.String()
}

func AddWriter(w io.Writer) {
	_logger.lock.Lock()
	defer _logger.lock.Unlock()

	_logger.loggers = append(_logger.loggers, w)
}

func Debug(facility Facility, message ...interface{}) {
	_logger.c <- &Message{facility, LevelDebug, _string(message)}
}

func Info(facility Facility, message ...interface{}) {
	_logger.c <- &Message{facility, LevelInfo, _string(message)}
}

func Warn(facility Facility, message ...interface{}) {
	_logger.c <- &Message{facility, LevelWarn, _string(message)}
}

func Error(facility Facility, message ...interface{}) {
	_logger.c <- &Message{facility, LevelError, _string(message)}
}

func Fatal(facility Facility, message ...interface{}) {
	_logger.c <- &Message{facility, LevelFail, _string(message)}
}

func Cancel() {
	_logger.c <- &Message{level: LevelCancel}
}

func Run(minLevel Level) bool {
	for m := range _logger.c {
		if m.level == LevelCancel {
			return true
		}

		if m.level < minLevel {
			continue
		}

		t := time.Now().Format(time.RFC3339)
		s := []byte(_string([]interface{}{t, m.level.String() + ":", m.facility + ":", m.message}))

		_logger.lock.Lock()

		for _, l := range _logger.loggers {
			l.Write(s)
		}

		_logger.lock.Unlock()

		if m.level == LevelFail {
			return false
		}
	}

	return true
}
