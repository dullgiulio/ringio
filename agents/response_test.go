package agents

import (
	"errors"
	"testing"
)

func TestResponseString(t *testing.T) {
	rs := NewAgentMessageResponseString()
	r := AgentMessageResponse(&rs)

	go func() {
		r.Err(errors.New("Some error"))
	}()

	d, e := r.Get()

	if d.(string) != "" {
		t.Error("Expected response data to be empty string")
	}

	if e == nil {
		t.Error("Expected response error")
	}

	if e.Error() != "Some error" {
		t.Error("Unexpected error returned")
	}

	rs = NewAgentMessageResponseString()
	r = AgentMessageResponse(&rs)

	go func() {
		r.Data("some string")
		r.Ok()
	}()

	d, e = r.Get()

	if d.(string) != "some string" {
		t.Error("Unexpected data in response")
	}

	if e != nil {
		t.Error("Unexpected error in response")
	}
}

func TestResponseInt(t *testing.T) {
	rs := NewAgentMessageResponseInt()
	r := AgentMessageResponse(rs)

	go func() {
		r.Err(errors.New("Some error"))
	}()

	d, e := r.Get()

	if d.(int) != 0 {
		t.Error("Expected response data to be empty string")
	}

	if e == nil {
		t.Error("Expected response error")
	}

	if e.Error() != "Some error" {
		t.Error("Unexpected error returned")
	}

	rs = NewAgentMessageResponseInt()
	r = AgentMessageResponse(rs)

	go func() {
		r.Data(100)
		r.Ok()
	}()

	d, e = r.Get()

	if d.(int) != 100 {
		t.Error("Unexpected data in response")
	}

	if e != nil {
		t.Error("Unexpected error in response")
	}
}

func TestResponseBool(t *testing.T) {
	rs := NewAgentMessageResponseBool()
	r := AgentMessageResponse(&rs)

	go func() {
		r.Err(errors.New("Some error"))
	}()

	d, e := r.Get()

	if d.(bool) != false {
		t.Error("Expected response data to be empty string")
	}

	if e == nil {
		t.Error("Expected response error")
	}

	if e.Error() != "Some error" {
		t.Error("Unexpected error returned")
	}

	rs = NewAgentMessageResponseBool()
	r = AgentMessageResponse(&rs)

	go func() {
		// Calling .Data() on Bool doesn't do anything by choice.
		r.Ok()
	}()

	d, e = r.Get()

	if d.(bool) != true {
		t.Error("Unexpected data in response")
	}

	if e != nil {
		t.Error("Unexpected error in response")
	}
}

func TestResponseList(t *testing.T) {
	rs := NewAgentMessageResponseList()
	r := AgentMessageResponse(&rs)

	go func() {
		r.Err(errors.New("Some error"))
	}()

	d, e := r.Get()

	if d.(*[]AgentDescr) != nil {
		t.Error("Expected response data to be nil")
	}

	if e == nil {
		t.Error("Expected response error")
	}

	if e.Error() != "Some error" {
		t.Error("Unexpected error returned")
	}

	rs = NewAgentMessageResponseList()
	r = AgentMessageResponse(&rs)

	go func() {
		ag := []AgentDescr{
			AgentDescr{},
			AgentDescr{},
		}

		r.Data(&ag)
		r.Ok()
	}()

	d, e = r.Get()

	ag := d.(*[]AgentDescr)
	if len(*ag) != 2 {
		t.Error("Unexpected data in response")
	}

	if e != nil {
		t.Error("Unexpected error in response")
	}
}
