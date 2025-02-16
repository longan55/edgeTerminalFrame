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
	version      = "V0.0.0209"
	GATEWAYMODE  = 1 //网关模式
	PLATFORMMODE = 0 //平台模式
)

type Task struct {
	F       func() error
	Content string
}

type process struct {
	mux      sync.Mutex
	runmode  int
	quittask []Task
}

func (p *process) SetRunMode(mode int) {
	p.runmode = mode
}
func (p *process) RunMode() int {
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
	Logger.Info("EDGE程序启动成功", zap.String("Version", version))
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
