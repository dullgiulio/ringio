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
	ac       *agents.Collection
	resp     agents.AgentMessageResponseBool
	actionCh chan ServerAction
	socket   string
	over     bool
	mux      *sync.Mutex
}

func NewRpcServer(socket string, autorun bool) *RpcServer {
	ac := agents.NewCollection()
	s := &RpcServer{
		socket:   socket,
		ac:       ac,
		resp:     agents.NewAgentMessageResponseBool(),
		actionCh: make(chan ServerAction),
		mux:      &sync.Mutex{},
	}

	go ac.Run(autorun)

	return s
}

type RpcReq struct {
	Action ServerAction
	Agent  *agents.AgentDescr
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

func Run(sessionName string) (returnStatus int) {
	socket := utils.FileInDotpath(sessionName)
	s := NewRpcServer(socket, config.C.AutoRun)

	// Serve RPC calls.
	go s.serve()

	// Relay errors to the main logging.
	go func() {
		s.relayErrorsAndOutput()

		log.Cancel()
	}()

	// Print all logged information.
	// TODO: minimum level is set via --flag.
	if !log.Run(log.LevelDebug) {
		returnStatus = 1
	}

	s.cleanup()

	return
}

func (s *RpcServer) Add(req *RpcReq, result *int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	ag := req.Agent
	resp := agents.NewAgentMessageResponseInt()

	switch ag.Type {
	case agents.AgentTypePipe:
		s.ac.Add(agents.NewAgentPipe(ag.Args[0], ag.Meta.Role, ag.Meta.Filter), resp)
	case agents.AgentTypeCmd:
		s.ac.Add(agents.NewAgentCmd(ag.Args, ag.Meta.Role, ag.Meta.Filter), resp)
	case agents.AgentTypeNull:
		s.ac.Add(agents.NewAgentNull(0, ag.Meta.Role, ag.Meta.Filter), resp)
	}

	i, err := resp.Get()
	*result = i.(int)

	return err
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

func (s *RpcServer) Stop(id int, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	na := agents.NewAgentNull(id, agents.AgentRoleSink, nil) // Role is unimportant here.
	s.ac.SetAgentStatusKill(na, &s.resp)
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

	onexit.PendingExit(0)

	return nil
}

func (s *RpcServer) Set(req *RpcReq, result *RpcResp) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.over {
		return sessionOver
	}

	*result = true

	s.actionCh <- req.Action

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
