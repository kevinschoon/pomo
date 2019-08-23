package main

import (
	"encoding/json"
	"net"
)

var _ Client = (*SocketClient)(nil)

// SocketClient makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type SocketClient struct {
	conn net.Conn
}

func (c SocketClient) read(statusCh chan *Status) {
	buf := make([]byte, 512)
	n, _ := c.conn.Read(buf)
	status := &Status{}
	json.Unmarshal(buf[0:n], status)
	statusCh <- status
}

func (c SocketClient) Status() (*Status, error) {
	statusCh := make(chan *Status)
	c.conn.Write([]byte("status"))
	go c.read(statusCh)
	return <-statusCh, nil
}

func (c SocketClient) Suspend() bool {
	panic("not implemented")
}

func (c SocketClient) Toggle() {
	panic("not implemented")
}

func (c SocketClient) Close() error { return c.conn.Close() }

func NewSocketClient(path string) (*SocketClient, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	return &SocketClient{conn: conn}, nil
}
