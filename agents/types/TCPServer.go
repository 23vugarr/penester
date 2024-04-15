package types

import (
	"fmt"
	"log"
	"net"
)

type TCPServer struct {
	IP       string
	Port     string
	Listener net.Listener
}

func NewTcpServer(ip, port string) *TCPServer {
	return &TCPServer{
		IP:   ip,
		Port: port,
	}
}

func (t *TCPServer) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s%s", t.IP, t.Port))
	if err != nil {
		return err
	}
	t.Listener = listener

	return nil
}

func (t *TCPServer) Close() error {
	err := t.Listener.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (t *TCPServer) Accept() (net.Conn, error) {
	conn, err := t.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
