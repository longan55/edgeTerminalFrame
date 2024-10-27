package core

type Connector interface {
	Connect() error
	Disconnect() error
	//查看启停状态
	OnOff() bool
	//设置启/停
	SetOnOff(bool)
	//Ping 查看在线状态
	Ping() bool
	//ChangeState 更新在线状态
	ChangeState(bool)
	//根据业务场景开启自动重连，不可关闭。
	AutoReconnect()
	//设置重连间隔
	SetReconnectCycle(cycle int)
	//设置重连超时时间
	SetReconnectTimeout(timeout int)
}

type ConnectMan struct {
	//启停状态
	onoff bool
	//连接状态
	online bool
	//是否自动重连
	autoReconnect bool
	//自动重连间隔,单位s
	reconnectCycle int
	//重连超时时间,单位s
	reconnectTimeout int
}

func (man *ConnectMan) Connect() error {
	return nil
}

func (man *ConnectMan) DisConnect() error {
	return nil
}
func (man *ConnectMan) OnOff() bool {
	return man.onoff
}

func (man *ConnectMan) SetOnOff(onoff bool) {
	man.onoff = onoff
}

func (man *ConnectMan) Ping() bool {
	return man.online
}

func (man *ConnectMan) ChangeState(online bool) {
	man.online = online
}

func (man *ConnectMan) AutoReconnect() {
	man.autoReconnect = true
}

func (man *ConnectMan) SetReconnectCycle(cycle int) {
	man.reconnectCycle = cycle
}

func (man *ConnectMan) SetReconnectTimeout(timeout int) {
	man.reconnectTimeout = timeout
}

type Sender interface {
	Send(data []byte) error
}

// BreakPointer 断点续传接口
type BreakPointer interface {
	Point()
}

type Point struct {
}
