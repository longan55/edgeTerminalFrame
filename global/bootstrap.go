package global

import (
	// "cnc-edge/global"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

var Process = &process{}

const (
	version              = "V0.0.0217"
	GATEWAYMODE  runmode = 0 //网关模式
	PLATFORMMODE runmode = 1 //平台模式
)

type runmode byte

func (r runmode) String() string {
	switch r {
	case 0:
		return "网关模式"
	case 1:
		return "平台模式"
	}
	return "未知"
}

type Task struct {
	F       func() error
	Content string
}

type process struct {
	mux      sync.Mutex
	runmode  runmode
	quittask []Task
}

func (p *process) SetRunMode(mode runmode) {
	p.runmode = mode
}
func (p *process) RunMode() runmode {
	return p.runmode
}

func (p *process) RegisterQuitTask(task Task) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.quittask = append(p.quittask, Task{
		task.F, task.Content,
	})
}

func RegisterQuitTask(task Task) {
	Process.mux.Lock()
	defer Process.mux.Unlock()

	Process.quittask = append(Process.quittask, Task{
		task.F, task.Content,
	})
}

func GracefullyExit() {
	Logger.Info("EDGE程序启动成功", zap.String("模式", Process.runmode.String()), zap.String("版本", version))
	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	sig := <-close
	Logger.Info("程序接收信号", zap.String("Signal", sig.String()), zap.Any("Code", fmt.Sprintf("%d", sig)))
	for _, task := range Process.quittask {
		Logger.Info("执行退出任务: " + task.Content)
		if err := task.F(); err != nil {
			Logger.Error(task.Content, zap.Error(err))
		}
	}
}
