package agents

type AgentMessageResponse interface {
	Ok()
	Err(error)
	Data(interface{})
	Get() (interface{}, error)
}

type AgentMessageResponseString struct {
	readyCh chan bool
	data    string
	err     error
}

func NewAgentMessageResponseString() AgentMessageResponseString {
	return AgentMessageResponseString{
		readyCh: make(chan bool),
	}
}

func (r *AgentMessageResponseString) Get() (interface{}, error) {
	<-r.readyCh
	return r.data, r.err
}

func (r *AgentMessageResponseString) Ok() {
	r.readyCh <- true
}

func (r *AgentMessageResponseString) Err(err error) {
	r.err = err
	r.readyCh <- false
}

func (r *AgentMessageResponseString) Data(data interface{}) {
	r.data = data.(string)
}

type AgentMessageResponseInt struct {
	i   int
	err chan error
}

func NewAgentMessageResponseInt() *AgentMessageResponseInt {
	return &AgentMessageResponseInt{
		err: make(chan error),
	}
}

func (r *AgentMessageResponseInt) Get() (i interface{}, e error) {
	return r.i, <-r.err
}

func (r *AgentMessageResponseInt) Ok() {
	r.err <- nil
}

func (r *AgentMessageResponseInt) Err(err error) {
	r.err <- err
}

func (r *AgentMessageResponseInt) Data(i interface{}) {
	r.i = i.(int)
}

type AgentMessageResponseBool chan error

func NewAgentMessageResponseBool() AgentMessageResponseBool {
	return make(chan error)
}

func (r AgentMessageResponseBool) Get() (interface{}, error) {
	err := <-r
	val := false

	if err == nil {
		val = true
	}

	return val, err
}

func (r AgentMessageResponseBool) Ok() {
	r <- nil
}

func (r AgentMessageResponseBool) Err(err error) {
	r <- err
}

func (r AgentMessageResponseBool) Data(interface{}) {
	// Do nothin'
}

type AgentMessageResponseList struct {
	errorCh chan error
	Agents  *[]AgentDescr
	Error   error
}

func NewAgentMessageResponseList() AgentMessageResponseList {
	return AgentMessageResponseList{
		errorCh: make(chan error),
	}
}

func (r *AgentMessageResponseList) Get() (interface{}, error) {
	r.Error = <-r.errorCh
	return r.Agents, r.Error
}

func (r *AgentMessageResponseList) Ok() {
	r.errorCh <- nil
}

func (r *AgentMessageResponseList) Err(err error) {
	r.errorCh <- err
}

func (r *AgentMessageResponseList) Data(data interface{}) {
	r.Agents = data.(*[]AgentDescr)
}
