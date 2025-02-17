package main

import (
	"edgeTerminalFrame/comm"
	"edgeTerminalFrame/core"
	"edgeTerminalFrame/global"
	"edgeTerminalFrame/gopool"
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"
)

func bootStrap() {
	global.LoadConfig(true)

	//init logger
	global.InitLogger()
	gopool.SetLogger(global.Logger)
	// 3.加载数据库
	global.InitBoltdb()

	// 4.连接池、管理重连和状态
	comm.InitConnect()
	global.Logger.Info("连接管理池启动成功")
}

func main() {
	bootStrap()
	defer func() {
		if e := recover(); e != nil {
			global.Logger.Panic(fmt.Sprintf("%+v\n%v", e, fmt.Sprint(string(debug.Stack()))))
		}
	}()

	gopool.Go(func() {
		if err := core.EdgeCore.Preload(); err != nil {
			global.Logger.Error("Preload", zap.Error(err))
		}
	})
	// 监听退出信号->优雅退出
	global.GracefullyExit()
}

// 1. 终端基础信息管理
// 2. 设备信息
// 3. 抽象设备管理, 结合几个实体设备. 绑定关系
// 4. 状态热更新,信息热更新
// 5. 重连管理
// 6. 启动/关闭管理
