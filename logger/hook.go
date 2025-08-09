package logger

import (
	"io"
	"path/filepath"

	"quickstart-go-react/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LevelFileHook 按级别分别记录到不同文件的钩子
type LevelFileHook struct {
	writers map[logrus.Level]io.Writer
}

// NewLevelFileHook 创建新的级别文件钩子
func NewLevelFileHook(logConfig config.LogConfig) *LevelFileHook {
	hook := &LevelFileHook{
		writers: make(map[logrus.Level]io.Writer),
	}

	// 为不同级别创建不同的日志文件
	levelFiles := map[logrus.Level]string{
		logrus.TraceLevel: "trace.log",
		logrus.DebugLevel: "debug.log",
		logrus.InfoLevel:  "info.log",
		logrus.WarnLevel:  "warn.log",
		logrus.ErrorLevel: "error.log",
		logrus.FatalLevel: "fatal.log",
		logrus.PanicLevel: "panic.log",
	}

	for level, filename := range levelFiles {
		hook.writers[level] = &lumberjack.Logger{
			Filename:   filepath.Join(logConfig.OutputDir, filename),
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}
	}

	return hook
}

// Levels 返回钩子处理的日志级别
func (hook *LevelFileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire 处理日志条目
func (hook *LevelFileHook) Fire(entry *logrus.Entry) error {
	if writer, ok := hook.writers[entry.Level]; ok {
		// 格式化日志条目
		formatter := &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}

		serialized, err := formatter.Format(entry)
		if err != nil {
			return err
		}

		_, err = writer.Write(serialized)
		return err
	}
	return nil
}
