package jlog

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"testing"
)

func Test普通日志(t *testing.T) {
	c := LogConfig{
		Level:         "info",
		Dir:           ".tmp/logs",
		FilePattern:   "%Y-%m-%d/joo-${level}.log", // 按日期和等级分割
		MaxAge:        0,
		Format:        "json",
		EncodeLevel:   "",
		ShowLine:      true,
		AppendConsole: true,
		LogKey: &KeyLogConfig{
			TimeKey:       "",
			LevelKey:      "",
			NameKey:       "",
			CallerKey:     "",
			FunctionKey:   "",
			MessageKey:    "",
			StacktraceKey: "",
		},
	}
	c.FixDefault()
	logger := Zap.New(c)
	fn := func(l *zap.Logger) {
		l.Debug("这是 debug")
		l.Info("这是 info ")
		l.Warn("这是 warn")
		l.Error("这是 error")
		l.DPanic("这是 dpanic")
		//l.Panic("这是 panic")
	}
	fn(logger)
}

// NOTE 序列化失败
func Test标准配置(t *testing.T) {
	fnDump := func(name string, v interface{}) (string, error) {
		var err error
		var b1 bytes.Buffer
		switch name {
		case "json":
			enc := json.NewEncoder(&b1)
			enc.SetIndent("", "  ")
			err = enc.Encode(v)
		case "yaml":
			err = yaml.NewEncoder(&b1).Encode(v)
		case "toml":
			err = toml.NewEncoder(&b1).Encode(v)
		}
		if err != nil {
			return "", err
		}
		return b1.String(), err
	}
	fnShow := func(e string, c zap.Config) {
		t.Logf("%s%#v", e, c)
		if s, err := fnDump("json", c); err != nil {
			t.Error(err)
		} else {
			t.Log(e, s)
		}
		if s, err := fnDump("yaml", c); err != nil {
			t.Error(err)
		} else {
			t.Log(e, s)
		}
		if s, err := fnDump("toml", c); err != nil {
			t.Error(err)
		} else {
			t.Log(e, s)
		}

	}

	cfgProd := zap.NewProductionConfig()
	fnShow("prod: ", cfgProd)

	cfgDev := zap.NewDevelopmentConfig()
	fnShow("dev: ", cfgDev)
}
