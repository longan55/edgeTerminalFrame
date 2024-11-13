package connect

import (
	"fmt"
	"net"
)

type Connector interface {
	Connect() error
	Close() error
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
