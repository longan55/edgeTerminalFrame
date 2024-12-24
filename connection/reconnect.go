package connect

import (
	"edgeTerminalFrame/global"
	"edgeTerminalFrame/gopool"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

func init() {
	ConnectorManager = &connectorManager{
		container: make(map[Connector]bool),
		mutex:     sync.RWMutex{},
	}
	gopool.Go(func() {
		ticker2 := time.NewTicker(12500 * time.Millisecond).C
		for range ticker2 {
			for c, isConnected := range ConnectorManager.container {
				if !isConnected {
					ConnectorManager.connect(c)
				}
			}
		}
	})
}

var ConnectorManager *connectorManager

type connectorManager struct {
	container map[Connector]bool
	mutex     sync.RWMutex
}

func (manager *connectorManager) AddConnector(c Connector, state bool) {
	//已添加
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.container[c] = state
}

func (manager *connectorManager) DelConnector(c Connector) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.container, c)
}

func (manager *connectorManager) State(c Connector) bool {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	return manager.container[c]
}

func (manager *connectorManager) connect(c Connector) {
	gopool.Go(func() {
		c.Close()
		err := c.Connect()
		//连接失败 -> 返回（状态本为false无需修改）
		if err != nil {
			global.Logger.Error("重连失败", zap.String("*Connector", fmt.Sprintf("%p", &c)), zap.String("Address", c.Uri()), zap.Error(err))
			return
		}
		//连接成功 -> 修改状态
		global.Logger.Info("重连成功", zap.String("Connector", c.Uri()))
		manager.AddConnector(c, true)
	})
}
