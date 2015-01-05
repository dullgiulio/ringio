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

type RpcServer struct {
	ac     *agents.Collection
	resp   agents.AgentMessageResponseBool
	socket string
	over   bool
	mux    *sync.Mutex
}

func NewRpcServer(socket string, autorun bool) *RpcServer {
	ac := agents.NewCollection()
	s := &RpcServer{
		socket: socket,
		ac:     ac,
		resp:   agents.NewAgentMessageResponseBool(),
		mux:    &sync.Mutex{},
	}

	go ac.Run(autorun)

	return s
}

type RpcReq struct {
	Agent *agents.AgentDescr
}

type RpcResp bool

var sessionOver error = errors.New("Session has already terminated")

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
	s := NewRpcServer(socket, config.C.AutoRun)

	// Serve RPC calls.
	go s.serve()

	// Print all logged information.
	if !log.Run(logLevel) {
		returnStatus = 1
	}

	s.cleanup()

	return
}

func (s *RpcServer) Ping(req *RpcReq, result *int) error {
	*result = 1

	// XXX: It would be nice to verify that the agents routine is actually running.

	return nil
}

func (s *RpcServer) Add(req *RpcReq, result *int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	var agent agents.Agent

	if s.over {
		return sessionOver
	}

	ag := req.Agent
	resp := agents.NewAgentMessageResponseInt()

	switch ag.Type {
	case agents.AgentTypePipe:
		agent = agents.NewAgentPipe(ag.Args[0], ag.Meta.Role, ag.Meta.Filter, ag.Meta.Options)
	case agents.AgentTypeCmd:
		agent = agents.NewAgentCmd(ag.Args, ag.Meta.Role, ag.Meta.Filter, ag.Meta.Options)
	case agents.AgentTypeNull:
		agent = agents.NewAgentNull(0, ag.Meta.Role, ag.Meta.Filter, ag.Meta.Options)
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

func (s *RpcServer) List(req *RpcReq, result *agents.AgentMessageResponseList) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	*result = agents.NewAgentMessageResponseList()

	s.ac.List(result)
	result.Get()

	return result.Error
}

func (s *RpcServer) Run(req *RpcReq, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	*result = true

	s.ac.StartAll(&s.resp)
	s.resp.Get()

	return nil
}

func (s *RpcServer) StartAll(req *RpcReq, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	s.ac.StartAll(&s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RpcServer) Start(id int, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	na := agents.NewAgentNull(id, agents.AgentRoleSink, nil, &agents.AgentOptions{}) // Role is unimportant here.
	s.ac.Start(na, &s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RpcServer) Stop(id int, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	na := agents.NewAgentNull(id, agents.AgentRoleSink, nil, &agents.AgentOptions{}) // Role is unimportant here.
	s.ac.Stop(na, &s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RpcServer) Kill(id int, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	na := agents.NewAgentNull(id, agents.AgentRoleSink, nil, &agents.AgentOptions{}) // Role is unimportant here.
	s.ac.Kill(na, &s.resp)
	s.resp.Get()

	*result = true
	return nil
}

func (s *RpcServer) Close(req *RpcReq, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
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

func (s *RpcServer) cleanup() {
	s.mux.Lock()
	defer s.mux.Unlock()

	// This session has terminated.
	s.over = true
	os.Remove(s.socket)
}

func (s *RpcServer) serve() {
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
