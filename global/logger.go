package global

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	rotate "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConf struct {
	Level        string `yaml:"level"`
	Path         string `yaml:"path"`
	Partten      string `yaml:"partten"`
	MaxAge       int    `yaml:"maxAge"`
	RotationTime int    `yaml:"rotationTime"`
	Compress     bool   `yaml:"compress"`
}

var Logger *zap.Logger

func InitLogger() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	logfile, err := os.OpenFile("./native.log", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	if err := initLogger(); err != nil {
		log.SetOutput(io.MultiWriter(logfile, os.Stdout))
		log.Printf("init logger error: %v\n", err)
	}
}

func initLogger() error {
	if err := os.MkdirAll(AppConfig.LogConf.Path, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("create log dir [%s] error: %w", AppConfig.LogConf.Path, err)
	}

	var writer *rotate.RotateLogs
	var err error

	switch runtime.GOOS {
	case "windows":
		writer, err = rotate.New(
			path.Join(AppConfig.LogConf.Path, AppConfig.LogConf.Partten),
			rotate.WithMaxAge(time.Duration(AppConfig.LogConf.MaxAge)*24*time.Hour),          //文件最大保存时间
			rotate.WithRotationTime(time.Duration(AppConfig.LogConf.RotationTime)*time.Hour), //日志切割时间间隔
			rotate.WithHandler(rotate.HandlerFunc(CompressLog)),                              //注册 日志切割时回调函数-压缩日志
		)
	case "linux":
		writer, err = rotate.New(
			path.Join(AppConfig.LogConf.Path, AppConfig.LogConf.Partten),
			rotate.WithLinkName("latest.log"),                                                // 创建一个软链接指向最新的日志文件
			rotate.WithMaxAge(time.Duration(AppConfig.LogConf.MaxAge)*24*time.Hour),          //文件最大保存时间
			rotate.WithRotationTime(time.Duration(AppConfig.LogConf.RotationTime)*time.Hour), //日志切割时间间隔
			rotate.WithHandler(rotate.HandlerFunc(CompressLog)),                              //注册 日志切割时回调函数-压缩日志
		)
	}
	if err != nil {
		log.Fatalln("rotate.New:", err)
	}

	// 创建一个WriteSyncer，可以是os.Stdout、os.Stderr等等
	var ws zapcore.WriteSyncer

	switch AppConfig.System.Environment {
	case "devlopment", "test":
		ws = zapcore.AddSync(io.MultiWriter(writer, os.Stdout))
	default:
		ws = zapcore.AddSync(writer)
	}

	// 配置日志级别
	levelConf := zap.NewAtomicLevel()
	level, err := zapcore.ParseLevel(AppConfig.LogConf.Level)
	if err != nil {
		log.Printf("parse log level error: %v\n", err)
		levelConf.SetLevel(zapcore.InfoLevel)
	} else {
		levelConf.SetLevel(zapcore.Level(level))
	}

	// 编码器配置
	var encoderConfig zapcore.EncoderConfig
	switch AppConfig.System.Environment {
	case "production":
		encoderConfig = zap.NewProductionEncoderConfig()
	default:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建Encoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 创建core
	core := zapcore.NewCore(encoder, ws, levelConf)

	// 创建logger
	logger := zap.New(core).WithOptions(zap.AddCaller())

	Logger = logger

	Logger.Info("日志记录器创建成功")
	Logger.Info("配置文件", zap.Any("Content", AppConfig))
	return nil
}

func CompressLog(event rotate.Event) {
	//判断是否开启压缩
	if !AppConfig.LogConf.Compress {
		return
	}
	//判断是否是 日志切割 事件
	if event.Type() != rotate.FileRotatedEventType {
		return
	}
	fileevent := event.(*rotate.FileRotatedEvent)
	//上一个日志文件
	prePath := fileevent.PreviousFile()
	outputFile := prePath + ".gz"
	//prePath := "./log/" + preFile
	if prePath == "" {
		return
	}
	// 打开源文件
	inFile, err := os.Open(prePath)
	if err != nil {
		Logger.Error("compress log error: open log file fail", zap.String("FilePath", prePath), zap.Error(err))
		return
	}
	defer inFile.Close()
	//
	outFile, err := os.Create(outputFile)
	if err != nil {
		Logger.Error("compress log error: create compress file fail", zap.String("FilePath", prePath), zap.Error(err))
		return
	}
	defer outFile.Close()

	// 创建gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	buf := make([]byte, 1024*1024) // 1 MB 每次读取的块大小
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			break
		}

		// 将读取的数据写入gzip writer
		_, err = gzipWriter.Write(buf[:n])
		if err != nil {
			return
		}
	}
	err = os.Remove(prePath)
	if err != nil {
		Logger.Error("日志移除失败", zap.String("FilePath", prePath), zap.Error(err))
		return
	}
}