package core

// 虚拟设备

// 实体设备需要对接协议
//  A    B     C     D    E
// 平台 <=> 虚拟设备 <=> 实体设备
// 	   MQTT		    TCP、
//     HTTP API     HTTP、
//                  MODBUS、
//                  MQTT
// 1. 和实体设备通信:给实体设备发送指令、处理指令响应
// 2. 和平台通信
// 3. 信息热更新,更新信息立即生效.
// 3. 状态一致性,启停状态,在线状态,工作状态.
//  C                    E         B
// 设备, ———— 协议通道  数据源, 数据服务器, 模板
//		 \___ 协议通道  数据源, 数据服务器, 模板

// type Identity interface{
// 	Key()string
// 	Name()string
// }

// type IDCard struct{
// 	key string
// 	name string
// }

// func (card *IDCard)Key()string{
// 	return card.key
// }

// func (card *IDCard)Name()string{
// 	return card.name
// }

type Device struct {
	Key  string
	Name string
}

func NewDevice() *Device {
	return &Device{}
}

func (dev *Device) GetKey() string {
	return dev.Key
}
func (dev *Device) SetKey(key string) {
	dev.Key = key
}

func (dev *Device) GetName() string {
	return dev.Name
}
func (dev *Device) SetName(name string) {
	dev.Name = name
}
