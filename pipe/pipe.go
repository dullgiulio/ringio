package pipe

import (
	"fmt"
	"os"
	"syscall"

	"bitbucket.org/dullgiulio/ringio/log"
)

type Pipe struct {
	file     *os.File
	filename string
	isOpen   bool
}

func New(filename string) *Pipe {
	return &Pipe{
		filename: filename,
	}
}

func (p *Pipe) String() string {
	return p.filename
}

func (p *Pipe) Remove() error {
	return os.Remove(p.filename)
}

func (p *Pipe) OpenWrite() bool {
	if err := p.OpenWriteErr(); err != nil {
		log.Error(log.FacilityPipe, err)
		return false
	}

	log.Debug(log.FacilityPipe, "Pipe was opened for writing at", p.filename)
	return true
}

func (p *Pipe) OpenWriteErr() error {
	// XXX: File creation exploits are possible here. Woooohooo.
	if err := syscall.Mknod(p.filename, syscall.S_IFIFO|0600, 0); err != nil {
		return fmt.Errorf("Creating pipe for write failed: %s", err)
	}

	file, err := os.OpenFile(p.filename, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("Opening pipe for writing failed: %s", err)
	} else {
		p.isOpen = true
		p.file = file
	}

	return nil
}

func (p *Pipe) Write(b []byte) (n int, err error) {
	return p.file.Write(b)
}

func (p *Pipe) OpenRead() bool {
	if err := p.OpenReadErr(); err != nil {
		log.Error(log.FacilityPipe, err)
		return false
	}

	log.Debug(log.FacilityPipe, "Pipe was opened for reading at", p.filename)
	return true
}

func (p *Pipe) OpenReadErr() error {
	if file, err := os.Open(p.filename); err != nil {
		return fmt.Errorf("Opening pipe for reading failed: %s", err)
	} else {
		p.isOpen = true
		p.file = file
	}

	return nil
}

func (p *Pipe) Read(b []byte) (n int, err error) {
	return p.file.Read(b)
}

func (p *Pipe) Close() (err error) {
	if p.isOpen {
		err = p.file.Close()
		p.isOpen = false
	}

	return
}
