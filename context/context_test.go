package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type Server struct {
	worker Worker
}

func NewServer(worker Worker) *Server {
	return &Server{worker: worker}
}

type Worker interface {
	Start()
	Cancel()
}

type SpyWorker struct {
	t         *testing.T
	cancelled bool
	started   bool
}

func (s *SpyWorker) Start() {
	s.started = true
}

func (s *SpyWorker) Cancel() {
	s.cancelled = true
}

func (s *SpyWorker) AssertStartCalled() {
	if !s.cancelled {
		s.t.Errorf("worker was not started")
	}
}

func (s *SpyWorker) AssertCancelCalled() {
	if !s.cancelled {
		s.t.Errorf("worker was not cancelled")
	}
}

func (s *Server) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s.worker.Start()

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprint(w, "Hello there")
	case <-ctx.Done():
		s.worker.Cancel()
	}
}

func TestCancellation(t *testing.T) {
	worker := &SpyWorker{t: t}
	svr := NewServer(worker)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	// get request context and change so we auto send a cancellation after 1 nano
	ctx := req.Context()
	ctx, cancel := context.WithCancel(ctx)
	req = req.WithContext(ctx)
	time.AfterFunc(1*time.Nanosecond, cancel)

	svr.ServerHTTP(res, req)

	worker.AssertCancelCalled()
}
