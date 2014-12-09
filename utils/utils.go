package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

var _dotfileDir string
var _letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")

func getBasepath(home string) string {
	if home == "" {
		return path.Join("/tmp", ".ringio-"+os.Getenv("USER"))
	} else {
		return path.Join(home, ".ringio")
	}
}

func getDotfileDir() string {
	if _dotfileDir != "" {
		return _dotfileDir
	}

	rand.Seed(time.Now().UTC().UnixNano())
	basepath := getBasepath(os.Getenv("HOME"))

	if os.MkdirAll(basepath, 0750) != nil {
		panic("Unable to create directory for tagfile!")
	}

	_dotfileDir = basepath
	return _dotfileDir
}

func randSeq(n int) string {
	b := make([]rune, n)
	l := len(_letters)

	for i := range b {
		b[i] = _letters[rand.Intn(l)]
	}

	return string(b)
}

func GetRandomDotfile() string {
	return FileInDotpath(randSeq(8))
}

func FileInDotpath(filename string) string {
	return path.Join(getDotfileDir(), path.Base(filename))
}

func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "ringio: %s\n", err)
	os.Exit(1)
}

func Error(err error) {
    fmt.Fprintf(os.Stderr, "ringio: %s\n", err)
}
