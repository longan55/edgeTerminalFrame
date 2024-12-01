package main

import (
	"edgeTerminalFrame/core"
	"edgeTerminalFrame/global"
	"fmt"
)

func main() {
	global.InitConfig()
	//init logger
	global.InitLogger()
	//创建Host
	hostinfo := core.NewHostInfo()
	hostinfo.SetName("网关主体")
	hostinfo.SetSN("SNNNNNNNN")
	hostinfo.SetPosition("深圳")
	hostinfo.SetDescription("注释内容")
	host := core.NewHost(hostinfo)
	//
	cpuinfo := host.GetCpuInfo()
	fmt.Println("cpuinfo", cpuinfo)
	meminfo := host.GetMemoryInfo()
	fmt.Println("meminfo", meminfo)
	diskinfo := host.GetDiskInfo("/")
	fmt.Println("diskinfo", diskinfo)
	info := host.Info()
	fmt.Println("hostinfo", info)
}

// 1. 终端基础信息管理
// 2. 设备信息
// 3. 抽象设备管理, 结合几个实体设备. 绑定关系
// 4. 状态热更新,信息热更新
// 5. 重连管理
// 6. 启动/关闭管理
