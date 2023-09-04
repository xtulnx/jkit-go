package jlog

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
	"testing"
)

func Test解析(t *testing.T) {
	cfg0 := LogConfig{
		Level:         "info",
		Dir:           ".tmp/logs",
		FilePattern:   "",
		MaxAge:        3,
		Format:        "json",
		EncodeLevel:   "",
		ShowLine:      true,
		AppendConsole: true,
		LogKey:        nil,
	}
	cfg0.FixDefault()

	tests := []struct {
		name   string
		config string
	}{
		{"json", `{
  "level": "info",
  "dir": ".tmp/logs",
  "filePattern": "j-%Y-%m-%d.log",
  "maxAge": 3,
  "format": "json",
  "showLine": true,
  "appendConsole": true
}`},
		{"yaml", `
level: info
dir: .tmp/logs
filePattern: j-%Y-%m-%d.log
maxAge: 3
format: json
showLine: true
appendConsole: true
`},
		{"toml", `
Level = "info"
dir = ".tmp/logs"
filePattern = "j-%Y-%m-%d.log"
maxAge = 3
format = "json"
showLine = true
appendConsole = true
`},
	}

	for _, tt := range tests {
		t.Run("marshal "+tt.name, func(t *testing.T) {
			var err error
			var b1 bytes.Buffer
			switch tt.name {
			case "json":
				enc := json.NewEncoder(&b1)
				enc.SetIndent("", "  ")
				err = enc.Encode(cfg0)
			case "yaml":
				err = yaml.NewEncoder(&b1).Encode(cfg0)
			case "toml":
				err = toml.NewEncoder(&b1).Encode(cfg0)
			default:
				t.Error("unknown config type: ", tt.name)
			}
			if err != nil {
				t.Error(err)
				return
			} else if b1.String() != tt.config {
				t.Logf("config not equal: %s", b1.String())
				return
			}
		})
	}

	for _, tt := range tests {
		t.Run("unmarshal "+tt.name, func(t *testing.T) {
			var err error
			cfg := new(LogConfig)
			switch tt.name {
			case "json":
				err = json.Unmarshal([]byte(tt.config), cfg)
			case "yaml":
				err = yaml.Unmarshal([]byte(tt.config), cfg)
			case "toml":
				err = toml.Unmarshal([]byte(tt.config), cfg)
			default:
				t.Error("unknown config type: ", tt.name)
			}

			if err != nil {
				t.Error(err)
			} else {
				if cfg0 != *cfg {
					t.Errorf("config not equal %#v", *cfg)
				}
			}
		})
	}
}
