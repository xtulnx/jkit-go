package debug_zap

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xtulnx/jkit-go/jgorm/builder"
	"github.com/xtulnx/jkit-go/jgorm/config"
	"github.com/xtulnx/jkit-go/jgorm/debug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

// 使用 zap 的日志作为输出

func init() {
	// 默认用全局的 zap.L()
	builder.SetDefaultLoggerBuilder(defaultLoggerBuilder)
}

func defaultLoggerBuilder(c config.GeneralDB) logger.Interface {
	return newLoggerInner(zap.L(), zapcore.InfoLevel).kitLoggerBuilder(c)
}

// SetDefaultZapLogger 替换自定义的 zap 日志
func SetDefaultZapLogger(l *zap.Logger, level zapcore.Level) {
	builder.SetDefaultLoggerBuilder(newLoggerInner(l, level).kitLoggerBuilder)
}

type tLogger struct {
	logger *zap.Logger
	level  zapcore.Level
}

// Printf 日志打印
func (t tLogger) Printf(s string, i ...interface{}) {
	if t.logger == nil {
		return
	}
	if !t.logger.Level().Enabled(t.level) {
		return
	}
	var s1 string
	if len(i) > 0 {
		s1 = fmt.Sprintf(s, i...)
	} else {
		s1 = s
	}
	t.logger.Log(t.level, s1, zap.Skip())
}

// 创建一个调试日志
func (t *tLogger) kitLoggerBuilder(c config.GeneralDB) logger.Interface {
	l2 := builder.ToLogLevel(c.LogMode)
	l3 := ToZapLevel(c.LogZap)
	var w logger.Writer = t
	if zapcore.InvalidLevel.Enabled(l3) || t.logger == nil {
		w = log.New(os.Stdout, "\r\n", log.LstdFlags)
	} else if l3 != t.level {
		w = newLoggerInner(t.logger, l3)
	}
	return debug.NewLogger(w,
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      l2,
			Colorful:      true,
		},
	)
}

func ToZapLevel(l string) zapcore.Level {
	var l0 = zapcore.InvalidLevel
	if l != "" {
		switch strings.ToLower(l) {
		case "":
			//
		case "0", "false", "n", "no", "-1":
			l0 = zapcore.InvalidLevel
		case "1", "true", "y", "yes", "ok":
			l0 = zapcore.InfoLevel
		case "debug":
			l0 = zapcore.DebugLevel
		case "info":
			l0 = zapcore.InfoLevel
		case "warn", "warning":
			l0 = zapcore.WarnLevel
		case "error", "err":
			l0 = zapcore.ErrorLevel
		case "dpanic":
			l0 = zapcore.DPanicLevel
		case "panic":
			l0 = zapcore.PanicLevel
		case "fatal":
			l0 = zapcore.FatalLevel
		}
	}
	return l0
}

func newLoggerInner(l *zap.Logger, level zapcore.Level) *tLogger {
	if level == 0 {
		level = zapcore.InfoLevel
	}
	return &tLogger{logger: l, level: level}
}

func NewWriter(l *zap.Logger, level zapcore.Level) logger.Writer {
	return newLoggerInner(l, level)
}

func NewWriter1(l *zap.Logger, level string) logger.Writer {
	var level2 zapcore.Level
	if level == "" {
		level2 = zapcore.InfoLevel
	} else {
		level2 = ToZapLevel(level)
	}
	return newLoggerInner(l, level2)
}
