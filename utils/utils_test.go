package utils

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestDotdir(t *testing.T) {
	tmphome := fmt.Sprintf("/tmp/%d", time.Now().Unix()/int64(time.Millisecond))
	os.Setenv("HOME", tmphome)

	dir := getDotfileDir()
	if !strings.HasPrefix(dir, tmphome) {
		t.Error("Dot file was not created under the temporary home directory")
	}

	basepath := getBasepath("")

	if !strings.HasPrefix(basepath, "/tmp/") {
		t.Error("Base path without HOME is not in /tmp")
	}

	file := GetRandomDotfile()

	if file == "" {
		t.Error("Did not expect the filename to be nil")
	}

	if !strings.HasPrefix(file, tmphome) {
		t.Error("Dot file was not created under the temporary home directory")
	}

	if FileInDotpath(path.Base(file)) != file {
		t.Error("Basename does not equal to the dotpath")
	}

	if err := os.RemoveAll(tmphome); err != nil {
		t.Error(err)
	}
}
