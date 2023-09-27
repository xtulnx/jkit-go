package jlog

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

// ZapEncodeLevel
// "color": 小写编码器带颜色
// "capital": 大写编码器
// "capitalColor": 大写编码器带颜色
// 小写编码器(默认)
func (c *LogConfig) ZapEncodeLevel() zapcore.LevelEncoder {
	var e zapcore.LevelEncoder
	if err := e.UnmarshalText([]byte(c.EncodeLevel)); err == nil {
		return e
	}
	return zapcore.LowercaseLevelEncoder
}

// TransportLevel 获取日志级别，默认是 info
func (c *LogConfig) TransportLevel() zapcore.Level {
	var l zapcore.Level
	if e := l.UnmarshalText([]byte(strings.ToLower(c.Level))); e == nil {
		return l
	}
	return zapcore.InfoLevel
}

// GetLogPath 获取日志文件路径，如果没有配置，则返回空。如果只配置了目录，则按日期拆分
func (c *LogConfig) GetLogPath() string {
	if c.Dir != "" && c.FilePattern != "" {
		if path.IsAbs(c.FilePattern) {
			return c.FilePattern
		}
		return path.Join(c.Dir, c.FilePattern)
	} else if c.Dir != "" {
		return path.Join(c.Dir, DefaultFilePattern)
	} else {
		return c.FilePattern
	}
}

var Zap = new(_zap)

type _zap struct{}

func (z _zap) CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.000"))
}

func (z _zap) GetEncoderConfig(c *LogConfig) zapcore.EncoderConfig {
	ec := zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "logger",
		CallerKey:           "caller",
		FunctionKey:         "func",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         c.ZapEncodeLevel(),
		EncodeTime:          z.CustomTimeEncoder,
		EncodeDuration:      zapcore.SecondsDurationEncoder,
		EncodeCaller:        zapcore.FullCallerEncoder,
		EncodeName:          nil,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "",
	}
	if c.LogKey != nil {
		if c.LogKey.MessageKey != "" {
			ec.MessageKey = c.LogKey.MessageKey
		}
		if c.LogKey.LevelKey != "" {
			ec.LevelKey = c.LogKey.LevelKey
		}
		if c.LogKey.TimeKey != "" {
			ec.TimeKey = c.LogKey.TimeKey
		}
		if c.LogKey.NameKey != "" {
			ec.NameKey = c.LogKey.NameKey
		}
		if c.LogKey.CallerKey != "" {
			ec.CallerKey = c.LogKey.CallerKey
		}
		if c.LogKey.FunctionKey != "" {
			ec.FunctionKey = c.LogKey.FunctionKey
		}
		if c.LogKey.StacktraceKey != "" {
			ec.StacktraceKey = c.LogKey.StacktraceKey
		}
	}
	return ec
}

func (z _zap) GetEncoder(c *LogConfig) zapcore.Encoder {
	ec := z.GetEncoderConfig(c)
	if strings.ToLower(c.Format) == "json" {
		return zapcore.NewJSONEncoder(ec)
	}
	return zapcore.NewConsoleEncoder(ec)
}

func (z _zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel: // 调试级别
	case zapcore.InfoLevel: // 日志级别
	case zapcore.WarnLevel: // 警告级别
	case zapcore.ErrorLevel: // 错误级别
	case zapcore.DPanicLevel: // dpanic级别
	case zapcore.PanicLevel: // panic级别
	case zapcore.FatalLevel: // 终止级别
	default:
		level = zapcore.InfoLevel
	}
	return func(l zapcore.Level) bool {
		return l == level
	}
}

func (z _zap) GetWriteSyncer(c *LogConfig, logPath string) (zapcore.WriteSyncer, error) {
	var ws []io.Writer
	if logPath != "" {
		fileWriter, err := rotatelogs.New(logPath,
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithMaxAge(time.Duration(c.MaxAge)*24*time.Hour), // 日志留存时间
			rotatelogs.WithRotationTime(time.Hour*24),
		)
		if err == nil {
			if c.AppendConsole {
				ws = append(ws, os.Stdout)
			}
			ws = append(ws, fileWriter)
		}
	}
	var writer zapcore.WriteSyncer
	if len(ws) == 0 {
		writer = zapcore.AddSync(os.Stdout)
	} else if len(ws) == 1 {
		writer = zapcore.AddSync(ws[0])
	} else {
		ws1 := make([]zapcore.WriteSyncer, 0, len(ws))
		for _, w := range ws {
			ws1 = append(ws1, zapcore.AddSync(w))
		}
		writer = zapcore.NewMultiWriteSyncer(ws1...)
	}
	return writer, nil
}

func (z _zap) GetZapCores(c *LogConfig) []zapcore.Core {
	level := c.TransportLevel()
	encoder := z.GetEncoder(c)
	logpath := c.GetLogPath()
	split4level := strings.Index(logpath, HolderLevel) >= 0
	if logpath != "" {
		if split4level {
			cores := make([]zapcore.Core, 0, 7)
			for ; level <= zapcore.FatalLevel; level++ {
				p1 := strings.ReplaceAll(logpath, HolderLevel, level.String())
				writer, err := z.GetWriteSyncer(c, p1)
				if err == nil {
					cores = append(cores, zapcore.NewCore(encoder, writer, z.GetLevelPriority(level)))
				}
			}
			return cores
		}
	}
	writer, err := z.GetWriteSyncer(c, logpath)
	if err == nil {
		return []zapcore.Core{zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(level))}
	}
	return nil
}

func (z _zap) PrepareDir(d string) (ok bool, err error) {
	fi, err := os.Stat(d)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		} else {
			return false, errors.New("存在同名文件")
		}
	} else if os.IsNotExist(err) {
		fmt.Printf("创建日志目录: %s\n", d)
		err = os.MkdirAll(d, os.ModePerm)
		return err == nil, err
	}
	return false, err
}

// New 创建一个日志器
func (z _zap) New(c LogConfig) *zap.Logger {
	if c.Dir != "" {
		ok, err := z.PrepareDir(c.Dir)
		if err != nil {
			fmt.Printf("准备日志目录[%s]失败: %s\n", c.Dir, err.Error())
		}
		if !ok {
			c.Dir = ""
		}
	}
	cores := z.GetZapCores(&c)
	logger := zap.New(zapcore.NewTee(cores...))
	if c.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

type Logger zap.Logger

func (L *Logger) ForTask(name string) *Logger {
	return (*Logger)((*zap.Logger)(L).With(zap.String("task", name)))
}
