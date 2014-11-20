package config

type _Config struct {
	RingbufSize int64
	MaxLineSize int64
	AutoExit    bool
	AutoLock    bool
}

var C *_Config

var defaults _Config = _Config{
	RingbufSize: 1024,
	MaxLineSize: 1024,
	AutoExit:    true,
	AutoLock:    false,
}

func init() {
	C = &defaults
}

// This module is not thread safe. After initialization from flags,
// only read access is done on the immutable configuration option.
