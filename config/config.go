package config

import (
	"github.com/dullgiulio/ringbuf"
)

type _Config struct {
	RingbufSize    int64
	RingbufLogSize int64
	MaxLineSize    int64
	AutoExit       bool
	AutoLock       bool
	AutoRun        bool
	PrintLog       bool
	logring        *ringbuf.Ringbuf
}

var C *_Config

var defaults _Config = _Config{
	RingbufSize:    1024,
	RingbufLogSize: 1024,
	MaxLineSize:    1024,
	AutoRun:        true,
	AutoExit:       false,
	AutoLock:       false,
	PrintLog:       false,
}

func init() {
	C = &defaults
}

func Init() {
	C.logring = ringbuf.NewRingbuf(C.RingbufLogSize)
	go C.logring.Run()
}

func GetLogRingbuf() *ringbuf.Ringbuf {
	return C.logring
}

// This module is not thread safe. After initialization from flags,
// only read access is done on the immutable configuration option.
