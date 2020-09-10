package pomo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

// Server listens on a Unix domain socket
// for Pomo status requests
type Server struct {
	listener net.Listener
	runner   *TaskRunner
	running  bool
}

func (s *Server) listen() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}
		buf := make([]byte, 512)
		// Ignore any content
		conn.Read(buf)
		raw, _ := json.Marshal(s.runner.Status())
		conn.Write(raw)
		conn.Close()
	}
}

func (s *Server) Start() {
	s.running = true
	go s.listen()
}

func (s *Server) Stop() {
	s.running = false
	s.listener.Close()
}

func NewServer(runner *TaskRunner, config *Config) (*Server, error) {
	//check if socket file exists
	if _, err := os.Stat(config.SocketPath); err == nil {
		_, err := net.Dial("unix", config.SocketPath)
		//if error then sock file was saved after crash
		if err != nil {
			os.Remove(config.SocketPath)
		} else {
			// another instance of pomo is running
			return nil, errors.New(fmt.Sprintf("Socket %s is already in use", config.SocketPath))
		}
	}
	listener, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		return nil, err
	}
	return &Server{listener: listener, runner: runner}, nil
}

// Client makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type Client struct {
	conn net.Conn
}

func (c Client) read(statusCh chan *Status) {
	buf := make([]byte, 512)
	n, _ := c.conn.Read(buf)
	status := &Status{}
	json.Unmarshal(buf[0:n], status)
	statusCh <- status
}

func (c Client) Status() (*Status, error) {
	statusCh := make(chan *Status)
	c.conn.Write([]byte("status"))
	go c.read(statusCh)
	return <-statusCh, nil
}

func (c Client) Close() error { return c.conn.Close() }

func NewClient(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}
