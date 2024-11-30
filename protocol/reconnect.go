package connect

import "time"

type ConnectorManager struct {
	container map[Connector]bool
}

func (mgr *ConnectorManager) Add(c Connector) {
	mgr.container[c] = false
}

func (mgr *ConnectorManager) Del(c Connector) {
	delete(mgr.container, c)
}

func (mgr *ConnectorManager) Start() {
	go func() {
		check := time.NewTicker(5 * time.Second).C
		for {
			select {
			case <-check:
			}
		}
	}()
}
