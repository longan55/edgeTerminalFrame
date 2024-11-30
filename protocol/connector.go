package connect

import (
	"fmt"
	"net"
)

// Connector 主要用来管理重连
type Connector interface {
	Connect() error
	Close() error
	// Read(b []byte) (n int, err error)
	// Write(b []byte) (n int, err error)
	//LocalAddr() Addr
	//RemoteAddr() Addr
	//SetDeadline(t time.Time) error
	//SetReadDeadline(t time.Time) error
	//SetWriteDeadline(t time.Time) error
}

type TcpConnector struct {
	Options *TcpOptions
	conn    net.Conn
}

func NewTcpConnector(options *TcpOptions) *TcpConnector {
	return &TcpConnector{
		Options: options,
	}
}

func (connector *TcpConnector) Connect() error {
	conn, err := net.Dial("tcp", connector.Options.Address())
	if err != nil {
		return err
	}
	connector.conn = conn
	return nil
}

func (connector *TcpConnector) Close() error {
	return connector.conn.Close()
}

func (connector *TcpConnector) Read(b []byte) (n int, err error) {
	return connector.conn.Read(b)
}

func (connector *TcpConnector) Write(b []byte) (n int, err error) {
	return connector.conn.Write(b)
}

type TcpOptions struct {
	Protocol string
	Ip       string
	Port     string
}

func (options TcpOptions) Address() string {
	return fmt.Sprintf("%s:%s", options.Ip, options.Port)
}

func NewTcpOptions(ip, port string) *TcpOptions {
	return &TcpOptions{
		Ip:   ip,
		Port: port,
	}
}
