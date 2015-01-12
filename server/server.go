package server

import (
	"errors"
	"net"
	"net/rpc"
	"os"
	"sync"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/config"
	"github.com/dullgiulio/ringio/log"
	"github.com/dullgiulio/ringio/onexit"
	"github.com/dullgiulio/ringio/utils"
)

type RPCServer struct {
	ac     *agents.Collection
	resp   agents.AgentMessageResponseBool
	socket string
	over   bool
	mux    *sync.Mutex
}

func NewRPCServer(socket string, autorun bool) *RPCServer {
	ac := agents.NewCollection()
	s := &RPCServer{
		socket: socket,
		ac:     ac,
		resp:   agents.NewAgentMessageResponseBool(),
		mux:    &sync.Mutex{},
	}

	go ac.Run(autorun)

	return s
}

type RPCReq struct {
	Agent *agents.AgentDescr
}

type RPCResp bool

var errSessionOver = errors.New("Session has already terminated")

func Init() {
	if config.C.PrintLog {
		nw := log.NewNewlineWriter(os.Stderr)
		log.AddWriter(nw)
	}

	lw := ringbuf.NewBytes(config.GetLogRingbuf())
	log.AddWriter(lw)
}

func Run(logLevel log.Level, sessionName string) (returnStatus int) {
	socket := utils.FileInDotpath(sessionName)
	s := NewRPCServer(socket, config.C.AutoRun)

	// Serve RPC calls.
	go s.serve()

	// Print all logged information.
	if !log.Run(logLevel) {
		returnStatus = 1
	}

	s.cleanup()

	return
}

func (s *RPCServer) Ping(req *RPCReq, result *int) error {
	*result = 1

	// XXX: It would be nice to verify that the agents routine is actually running.

	return nil
}

func (s *RPCServer) Add(req *RPCReq, result *int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	var agent agents.Agent

	if s.over {
		return errSessionOver
	}

	ag := req.Agent
	resp := agents.NewAgentMessageResponseInt()

	switch ag.Type {
	case agents.AgentTypePipe:
		agent = agents.NewAgentPipe(ag.Args[0], &ag.Meta)
	case agents.AgentTypeCmd:
		agent = agents.NewAgentCmd(ag.Args, &ag.Meta)
	case agents.AgentTypeNull:
		agent = agents.NewAgentNull(0, &ag.Meta)
	}

	s.ac.Add(agent, resp)

	i, err := resp.Get()
	if err != nil {
		return err
	}

	*result = i.(int)

	if config.C.AutoRun || ag.Type == agents.AgentTypePipe {
		s.ac.Start(agent, &s.resp)

		if _, err := s.resp.Get(); err != nil {
			return err
		}
	}

	return nil
}

func (s *RPCServer) List(req *RPCReq, result *agents.AgentMessageResponseList) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return errSessionOver
	}

	*result = agents.NewAgentMessageResponseList()

	s.ac.List(result)
	result.Get()

	return result.Error
}

func (s *RPCServer) Run(req *RPCReq, result *RPCResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return errSessionOver
	}

	*result = true

	s.ac.StartAll(&s.resp)
	s.resp.Get()

	return nil
}

func (s *RPCServer) StartAll(req *RPCReq, result *RPCResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return errSessionOver
	}

	s.ac.StartAll(&s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RPCServer) Start(id int, result *RPCResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return errSessionOver
	}

	na := agents.NewAgentNull(id, agents.NewAgentMetadata()) // Role is unimportant here.
	s.ac.Start(na, &s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RPCServer) Stop(id int, result *RPCResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return errSessionOver
	}

	na := agents.NewAgentNull(id, agents.NewAgentMetadata()) // Role is unimportant here.
	s.ac.Stop(na, &s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RPCServer) Kill(id int, result *RPCResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return errSessionOver
	}

	na := agents.NewAgentNull(id, agents.NewAgentMetadata()) // Role is unimportant here.
	s.ac.Kill(na, &s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RPCServer) Close(req *RPCReq, result *RPCResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return errSessionOver
	}

	// This session has terminated.
	s.over = true

	*result = true

	s.ac.Cancel(&s.resp)
	s.resp.Get()

	log.Cancel()

	onexit.PendingExit(0)

	return nil
}

func (s *RPCServer) cleanup() {
	s.mux.Lock()
	defer s.mux.Unlock()

	// This session has terminated.
	s.over = true
	os.Remove(s.socket)
}

func (s *RPCServer) serve() {
	l, err := net.Listen("unix", s.socket)
	if err != nil {
		utils.Fatal(err)
	}

	onexit.Defer(s.cleanup)

	srv := rpc.NewServer()
	if err = srv.Register(s); err != nil {
		log.Fatal(log.FacilityDefault, err)
		onexit.PendingExit(1)
		return
	}

	srv.Accept(l)
}
