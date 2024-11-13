package core

import (
	"crypto/md5"
	"edgeTerminalFrame/gopool"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

// Host 宿主设备，一般是程序运行的设备，如网关、上位机...
type Host struct {
	info HostInfo
}

func NewHost(info HostInfo) *Host {
	return &Host{info: info}
}

func (host Host) Heartbeat() {

}

func (host Host) Name() string {
	return host.info.name
}

func (host Host) Info() HostInfo {
	return host.info
}

type HostInfo struct {
	name        string //名称
	sn          string //序列号
	position    string //地点
	longitude   string //经度
	latitude    string //纬度
	description string //描述
}

func NewHostInfo() HostInfo {
	return HostInfo{}
}
func (info *HostInfo) SetName(name string) {
	info.name = name
}
func (info *HostInfo) GetName() string {
	return info.name
}
func (info *HostInfo) SetSN(sn string) {
	info.sn = sn
}
func (info *HostInfo) GetSN() string {
	return info.sn
}
func (info *HostInfo) SetPosition(position string) {
	info.position = position
}
func (info *HostInfo) GetPosition() string {
	return info.position
}
func (info *HostInfo) SetLongitude(longitude string) {
	info.longitude = longitude
}
func (info *HostInfo) GetLongitude() string {
	return info.longitude
}
func (info *HostInfo) SetLatitude(latitude string) {
	info.latitude = latitude
}
func (info *HostInfo) GetLatitude() string {
	return info.latitude
}
func (info *HostInfo) SetDescription(description string) {
	info.description = description
}
func (info *HostInfo) GetDescription() string {
	return info.description
}

type CpuInfo struct {
	CpuNum    int     `json:"cpuNum"`    //cpu数量
	TotalUsed float64 `json:"totalUsed"` //CPU总使用率
	SelfUsed  float64 `json:"selfUsed"`  //当前进程CPU使用率
	Free      float64 `json:"free"`      //CPU空闲率
}

func (cpu CpuInfo) String() string {
	return fmt.Sprintf("CPU数量[%d] 已用[%.2f%%] 网关使用[%.2f%%] 空闲[%.2f%%]", cpu.CpuNum, cpu.TotalUsed, cpu.SelfUsed, cpu.Free)
}

// GetCpuInfo 获取cpu信息
func (host Host) GetCpuInfo() CpuInfo {
	var cpuInfo CpuInfo
	//cpu核心数
	counts, _ := cpu.Counts(true)
	cpuInfo.CpuNum = counts

	duration := 500 * time.Millisecond

	selfPercentChan := make(chan float64)
	cpuPercentChan := make(chan []float64)

	gopool.Go(func() {
		p, _ := process.NewProcess(int32(os.Getpid()))
		percent, err := p.Percent(duration)
		if err != nil {
			return
		}
		selfPercentChan <- percent
	})
	gopool.Go(func() {
		cpuPercent, err := cpu.Percent(duration, true)
		if err != nil {
			return
		}
		cpuPercentChan <- cpuPercent
	})

	selfPercent := <-selfPercentChan
	cpuPercent := <-cpuPercentChan
	//总CPU使用率
	var totalUsed float64
	for _, one := range cpuPercent {
		totalUsed += one
	}
	cpuInfo.TotalUsed = SaveBit(totalUsed, 2)
	//CPU空闲率
	cpuInfo.Free = 100 - cpuInfo.TotalUsed
	cpuInfo.Free = SaveBit(cpuInfo.Free, 2)
	//当前进程CPU使用率
	cpuInfo.SelfUsed = SaveBit(selfPercent, 2)

	return cpuInfo
}

type MemoryInfo struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
	Usage float64 `json:"usage"`
}

func (mem MemoryInfo) String() string {
	return fmt.Sprintf("总内存[%.2fG] 已用[%.2fG %.2f%%] 空闲[%.2fG]", mem.Total, mem.Used, mem.Usage, mem.Free)
}

func (host Host) GetMemoryInfo() MemoryInfo {
	var memory MemoryInfo
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return memory
	}
	//内存总量，GByte
	total := float64(vmStat.Total) / (1024 * 1024 * 1024)
	memory.Total = SaveBit(total, 2)

	//内存使用量
	used := float64(vmStat.Used) / (1024 * 1024 * 1024)
	memory.Used = SaveBit(used, 2)

	//内存空闲量
	free := float64(vmStat.Available) / (1024 * 1024 * 1024)
	memory.Free = SaveBit(free, 2)

	//使用率
	memory.Usage = SaveBit(vmStat.UsedPercent, 2)
	return memory
}

type DiskInfo struct {
	DirName string `json:"dirName"`
	//SysTypeName string  `json:"sysTypeName"`
	//TypeName    string  `json:"typeName"`
	Total float64 `json:"total"`
	Free  float64 `json:"free"`
	Used  float64 `json:"used"`
	Usage float64 `json:"usage"`
}

func (disk DiskInfo) String() string {
	return fmt.Sprintf("路径[%s] 总存储[%.2fG] 已用[%.2fG %.2f%%] 空闲[%.2fG]", disk.DirName, disk.Total, disk.Used, disk.Usage, disk.Free)
}

func (host Host) GetDiskInfo(path string) DiskInfo {
	var info DiskInfo
	info.DirName = path

	diskStat, err := disk.Usage(path)
	if err != nil {
		return info
	}
	//总量GB
	info.Total = float64(diskStat.Total) / (1024 * 1024 * 1024)
	//info.Total = fmt.Sprintf("%.02f GB", total)
	//空闲
	info.Free = float64(diskStat.Free) / (1024 * 1024 * 1024)
	//info.Free = fmt.Sprintf("%.02f GB", free)
	//已用
	info.Used = float64(diskStat.Used) / (1024 * 1024 * 1024)
	//info.Used = fmt.Sprintf("%.02f GB", used)
	//使用率
	info.Usage = SaveBit(diskStat.UsedPercent, 2)
	return info
}

// GetSerialNumber 获取网关序列号
func GetSerialNumber() (string, error) {
	addrs, err := MacAddress()
	if err != nil {
		return "", fmt.Errorf("获取mac地址失败:%w", err)
	}
	if len(addrs) == 0 {
		return "", errors.New("mac地址为空")
	}
	sum := md5.Sum([]byte(addrs[0]))

	sn := SerialNumberPrefix + fmt.Sprintf("%X", sum)[:8]
	return sn, nil
}

func MacAddress() (addrs []string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, inter := range interfaces {
		if inter.Flags&net.FlagLoopback == 0 {
			addr := inter.HardwareAddr.String()
			if addr != "" {
				addrs = append(addrs, strings.ToUpper(addr))
			}
		}
	}
	return
}
