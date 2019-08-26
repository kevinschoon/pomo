package main

import (
	"encoding/json"
	"net"
)

var _ Server = (*SocketServer)(nil)

// implement client for the primary usecase
// of running locally.
var _ Client = (*SocketServer)(nil)

// SocketServer listens on a Unix domain socket
// for Pomo status requests
type SocketServer struct {
	listener net.Listener
	task     *Task
	store    Store
	runner   *TaskRunner
	stop     chan struct{}
}

func (s *SocketServer) Serve() error {
	done := s.runner.Start()
	go func() {
		for {
			conn, err := s.listener.Accept()
			// TODO: fairly sure there is
			// a better way to handle this.
			if err != nil {
				return
			}
			buf := make([]byte, 512)
			// Ignore any content
			conn.Read(buf)
			status, _ := s.Status()
			raw, _ := json.Marshal(status)
			conn.Write(raw)
			conn.Close()
		}
	}()
loop:
	for {
		select {
		case <-done:
			break loop
		case <-s.stop:
			s.stop <- struct{}{}
			break loop
		}
	}
	err := s.store.With(func(st Store) error {
		return st.UpdateTask(s.task)
	})
	if err != nil {
		return err
	}
	return s.listener.Close()
}

func (s *SocketServer) Status() (*Status, error) {
	count := s.runner.Count()
	state := s.runner.State()
	timer := s.runner.Timer(count)
	return &Status{
		State:         state,
		Count:         count,
		Duration:      s.task.Duration,
		Message:       s.task.Message,
		NPomodoros:    len(s.task.Pomodoros),
		TimeStarted:   timer.TimeStarted(),
		TimeRunning:   timer.TimeRunning(),
		TimeSuspended: timer.TimeSuspended(),
	}, nil
}

func (s *SocketServer) Suspend() bool {
	return s.runner.Suspend()
}

func (s *SocketServer) Toggle() {
	s.runner.Toggle()
}

func (s *SocketServer) Stop() {
	s.stop <- struct{}{}
	<-s.stop
}

func NewSocketServer(task *Task, store Store, config *Config) (*SocketServer, error) {
	listener, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		return nil, err
	}
	return &SocketServer{
		listener: listener,
		store:    store,
		task:     task,
		runner:   NewTaskRunner(task),
		stop:     make(chan struct{}),
	}, nil
}
