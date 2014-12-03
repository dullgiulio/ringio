package server

import (
	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/log"
)

const (
	ActionNoop = iota
	ActionSetVerbose
	ActionUnsetVerbose
	ActionSetLocked
)

type ServerAction int

func (s *RpcServer) relayErrorsAndOutput() (returnStatus int) {
	var stdoutCh <-chan interface{}
	var verbose bool

	for {
		select {
		case data := <-stdoutCh:
			if data == nil {
				log.Debug(log.FacilityDefault, "Logging of stdout is complete")
				return
			}

			if verbose {
				log.Info(log.FacilityStdout, data)
			}
		case c := <-s.actionCh:
			switch c {
			case ActionSetVerbose:
				if !verbose {
					stdoutCh = s.ac.Reader(agents.CollectionRingTypeStdout)
				}

				verbose = true

				log.Info(log.FacilityDefault, "Set verbose mode on")
			case ActionUnsetVerbose:
				verbose = false

				log.Info(log.FacilityDefault, "Set verbose mode off")
			case ActionSetLocked:
				s.ac.Done()
			}
		}
	}
}
