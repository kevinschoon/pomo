package server

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/kevinschoon/pomo/pkg/runner"
)

var _ Server = (*SocketServer)(nil)

// SocketServer listens on a Unix domain socket
// for Pomo status requests
type SocketServer struct {
	socketPath string
	listener   net.Listener
	status     runner.Status
}

// NewSocketServer returns a new SocketServer
func NewSocketServer(socketPath string) *SocketServer {
	return &SocketServer{
		socketPath: socketPath,
	}
}

// SetStatus sets the most recent runner Status
func (s *SocketServer) SetStatus(status runner.Status) error {
	s.status = status
	return nil
}

// Start launches the SocketServer and listens for
// connections at the configured path
func (s *SocketServer) Start() error {
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	s.listener = listener
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			return err
		}
		buf := make([]byte, 512)
		// Ignore any content
		conn.Read(buf)
		raw, _ := json.Marshal(s.status)
		conn.Write(raw)
		conn.Close()
	}
}

// Stop stops the SocketServer
func (s *SocketServer) Stop() error {
	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		return err
	}
	return nil
}
