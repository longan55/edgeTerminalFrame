package comm

import (
	"edgeTerminalFrame/global"
	"errors"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

// 初始化
var (
	MqttCli *MqttClient
	once    sync.Once
)

type MQTTOptions struct {
	Broker string
	User   string
	Pass   string
	Topic  string
	//
	State bool
}

func InitMqtt(ops MQTTOptions) (err error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(ops.Broker)
	options.SetUsername(ops.User)
	options.SetPassword(ops.Pass)
	options.SetClientID(viper.GetString(global.CORE_SN))
	options.SetOnConnectHandler(onConnect)
	options.SetConnectionLostHandler(onConnectLost)
	options.SetDefaultPublishHandler(mqttHandler)
	//
	once.Do(func() {
		MqttCli = NewMqttSender(options)
		//注册退出任务: 断开Mqtt
		global.RegisterQuitTask(global.Task{
			F:       MqttCli.Close,
			Content: "断开Mqtt连接",
		})
	})
	//更新如何

	//连接
	if err = MqttCli.Connect(); err != nil {
		return
	}
	global.Logger.Info("mqtt broker: " + ops.Broker)
	global.Logger.Info("mqtt user: " + ops.User)
	global.Logger.Info("mqtt pass: " + ops.Pass)

	return nil
}

func mqttHandler(client mqtt.Client, message mqtt.Message) {

}

func onConnect(client mqtt.Client) {
	ConnectorManager.AddConnector(MqttCli, true)
	global.Logger.Info("mqtt Connected")
}
func onConnectLost(client mqtt.Client, err error) {
	// 自定义重连逻辑 todo 了解该库的Connect()是否自动重连
	ConnectorManager.AddConnector(MqttCli, false)
	global.Logger.Info("mqtt lost connection")
}

//

type MqttClient struct {
	options *mqtt.ClientOptions
	client  mqtt.Client
	mux     sync.Mutex
}

func NewMqttSender(options *mqtt.ClientOptions) *MqttClient {
	return &MqttClient{
		options: options,
	}
}

func (m *MqttClient) Connect() (err error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	//不再重连
	ConnectorManager.DelConnector(m)
	if m.client != nil {
		//更新mqtt连接信息时，需要先断连原来的连接，再进行连接
		m.client.Disconnect(50)
	}

	client := mqtt.NewClient(m.options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		return
	}

	m.client = client
	ConnectorManager.AddConnector(m, true)
	return nil
}

func (m *MqttClient) Close() error {
	m.mux.Lock()
	defer m.mux.Unlock()
	//不再重连
	ConnectorManager.DelConnector(m)
	//断连
	if m.client != nil {
		m.client.Disconnect(50)
	}
	return nil
}

// 发送数据
func (m *MqttClient) Send(topic string, data []byte) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if !ConnectorManager.State(m) {
		return errors.New("mqtt connection state is false")
	}

	if m.client == nil {
		return errors.New("MqttClient -> mqtt.Client is nil pointer")
	}
	token := m.client.Publish(topic, 0, false, data)
	return token.Error()
}

func (m *MqttClient) Uri() string {
	return viper.GetString(global.MQTT_BROKER)
}

func (m *MqttClient) Address() string {
	return ""
}

// 发布
func Publish(topic string, data []byte) error {
	if MqttCli != nil && MqttCli.client != nil {
		token := MqttCli.client.Publish(topic, 0, false, data)
		return token.Error()
	}
	return errors.New("nil pointer")
}

// 订阅
func Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if MqttCli != nil && MqttCli.client != nil {
		token := MqttCli.client.Subscribe(topic, qos, callback)
		return token.Error()
	}
	return errors.New("nil pointer")
}
