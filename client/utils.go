package client

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dullgiulio/ringio/pipe"
	"github.com/dullgiulio/ringio/utils"
)

func _removePipe(p *pipe.Pipe) {
	if err := p.Remove(); err != nil {
		utils.Error(err)
	}
}

func getAllSessions() []string {
	dir := utils.GetDotfileDir()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		utils.Error(err)
	}

	var sockets []string

	for _, finfo := range files {
		mode := finfo.Mode()

		if mode&os.ModeSocket == os.ModeSocket {
			sockets = append(sockets, finfo.Name())
		}
	}

	return sockets
}

func printList(list []string) {
	for _, l := range list {
		fmt.Printf("%s\n", l)
	}
}
