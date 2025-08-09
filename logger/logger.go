package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"quickstart-go-react/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	initialized bool
	initMutex   sync.Mutex
)

// Init 初始化日志系统
func Init() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	cfg := config.GetConfig()
	logConfig := cfg.Log

	// 创建日志目录
	if err := os.MkdirAll(logConfig.OutputDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 设置日志级别
	level, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		return fmt.Errorf("解析日志级别失败: %w", err)
	}
	logrus.SetLevel(level)

	// 设置日志格式
	if logConfig.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 配置日志输出
	var writers []io.Writer

	// app.log - 记录所有日志
	appLogWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logConfig.OutputDir, "app.log"),
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		Compress:   logConfig.Compress,
	}
	writers = append(writers, appLogWriter)

	// 如果启用控制台输出
	if logConfig.ConsoleOutput {
		writers = append(writers, os.Stdout)
	}

	// 设置多重输出
	multiWriter := io.MultiWriter(writers...)
	logrus.SetOutput(multiWriter)

	// 添加钩子，用于按级别分别记录到不同文件
	logrus.AddHook(NewLevelFileHook(logConfig))

	initialized = true
	logrus.Infof("日志系统初始化成功，日志目录: %s", logConfig.OutputDir)
	return nil
}

// ensureInit 确保日志系统已初始化
func ensureInit() {
	if !initialized {
		if err := Init(); err != nil {
			// 如果初始化失败，至少保证基本的日志功能
			logrus.SetLevel(logrus.InfoLevel)
			logrus.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
			})
			logrus.Errorf("日志系统初始化失败: %v", err)
		}
	}
}

// GetLogger 获取日志实例（确保初始化）
func GetLogger() *logrus.Logger {
	ensureInit()
	return logrus.StandardLogger()
}

func Traceln(args ...interface{}) {
	ensureInit()
	logrus.Traceln(args...)
}

func Debugln(args ...interface{}) {
	ensureInit()
	logrus.Debugln(args...)
}

func Infoln(args ...interface{}) {
	ensureInit()
	logrus.Infoln(args...)
}

func Warnln(args ...interface{}) {
	ensureInit()
	logrus.Warnln(args...)
}

func Errorln(args ...interface{}) {
	ensureInit()
	logrus.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	ensureInit()
	logrus.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	ensureInit()
	logrus.Panicln(args...)
}

func Tracef(format string, args ...interface{}) {
	ensureInit()
	logrus.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	ensureInit()
	logrus.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	ensureInit()
	logrus.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	ensureInit()
	logrus.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	ensureInit()
	logrus.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	ensureInit()
	logrus.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	ensureInit()
	logrus.Panicf(format, args...)
}
