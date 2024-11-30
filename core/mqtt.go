package core

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 初始化

func InitMqtt() error {
	options := mqtt.NewClientOptions()
	options.AddBroker("")
	options.SetUsername("")
	options.SetPassword("")
	options.SetOnConnectHandler(onConnect)
	options.SetConnectionLostHandler(onConnectLost)
	options.SetDefaultPublishHandler(mqttHandler)

	NewMqttSender(options)
	return nil
}

func mqttHandler(client mqtt.Client, message mqtt.Message) {

}

func onConnect(client mqtt.Client) {
	//global.Logger.Info("mqtt Connected")
}
func onConnectLost(client mqtt.Client, err error) {
	// 自定义重连逻辑 todo 了解该库的Connect()是否自动重连
	for i := 0; i < 3; i++ {
		// 尝试重新连接
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			break
		} else {
			time.Sleep(3)
		}
	}
}

//

type MqttClient struct {
	options *mqtt.ClientOptions
	client  mqtt.Client
}

func NewMqttSender(options *mqtt.ClientOptions) *MqttClient {
	return &MqttClient{
		options: options,
	}
}

func (m *MqttClient) Connect() (err error) {
	client := mqtt.NewClient(m.options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		return
	}
	m.client = client
	return nil
}

func (m *MqttClient) Disconnect() error {
	m.client.Disconnect(0)
	return nil
}

// 发送数据
func (m *MqttClient) Send(topic string, data []byte) error {
	token := m.client.Publish(topic, 0, false, data)
	return token.Error()
}

func (m *MqttClient) Uri() string {
	return ""
}

func (m *MqttClient) Address() string {
	return ""
}
