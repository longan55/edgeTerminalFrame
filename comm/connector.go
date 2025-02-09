package comm

// import (
// 	"fmt"
// 	"net"
// 	"time"
// )

// // Connector 主要用来管理重连
// type Connector interface {
// 	Connect() error
// 	Close() error
// 	Uri() string
// 	Address() string
// }

// type TcpConnector struct {
// 	Options *TcpOptions
// 	conn    net.Conn
// }

// func NewTcpConnector(options *TcpOptions) *TcpConnector {
// 	return &TcpConnector{
// 		Options: options,
// 	}
// }

// func (connector *TcpConnector) Connect() error {
// 	conn, err := net.Dial("tcp", connector.Address())
// 	if err != nil {
// 		return err
// 	}
// 	connector.conn = conn
// 	return nil
// }

// func (connector *TcpConnector) Close() error {
// 	return connector.conn.Close()
// }

// func (connector *TcpConnector) Uri() string {
// 	return connector.Address()
// }

// func (connector *TcpConnector) Address() string {
// 	return fmt.Sprintf("%s:%s", connector.Options.Ip, connector.Options.Port)
// }

// func (connector *TcpConnector) Read(b []byte) (n int, err error) {
// 	return connector.conn.Read(b)
// }

// func (connector *TcpConnector) Write(b []byte) (n int, err error) {
// 	return connector.conn.Write(b)
// }
// func (connector *TcpConnector) SetDeadline(t time.Time) error {
// 	return connector.conn.SetDeadline(t)
// }
// func (connector *TcpConnector) SetReadDeadline(t time.Time) error {
// 	return connector.conn.SetReadDeadline(t)
// }
// func (connector *TcpConnector) SetWriteDeadline(t time.Time) error {
// 	return connector.conn.SetWriteDeadline(t)
// }

// type TcpOptions struct {
// 	Protocol string
// 	Ip       string
// 	Port     string
// }

// func NewTcpOptions(ip, port string) *TcpOptions {
// 	return &TcpOptions{
// 		Ip:   ip,
// 		Port: port,
// 	}
// }
