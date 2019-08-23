package main

import (
	"encoding/json"
	"net"
)

var _ Server = (*SocketServer)(nil)

// SocketServer listens on a Unix domain socket
// for Pomo status requests
type SocketServer struct {
	listener net.Listener
	runner   *TaskRunner
}

func (s *SocketServer) Serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		buf := make([]byte, 512)
		// Ignore any content
		conn.Read(buf)
		status := &Status{
			State: s.runner.State(),
			Count: s.runner.Count(),
		}
		raw, _ := json.Marshal(status)
		conn.Write(raw)
		conn.Close()
	}
}

func NewSocketServer(runner *TaskRunner, config *Config) (*SocketServer, error) {
	listener, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		return nil, err
	}
	return &SocketServer{listener: listener, runner: runner}, nil
}
