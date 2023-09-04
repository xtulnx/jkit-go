package jlog

import "strings"

type KeyLogConfig struct {
	TimeKey       string
	LevelKey      string
	NameKey       string
	CallerKey     string
	FunctionKey   string
	MessageKey    string
	StacktraceKey string
}

const (
	HolderLevel        = "${level}"
	DefaultFilePattern = "j-%Y-%m-%d.log"
)

type LogConfig struct {
	// 级别
	Level string `mapstructure:"level,omitempty" json:"level,omitempty" yaml:"level,omitempty"`

	// 日志文件夹
	Dir string `mapstructure:"dir,omitempty" json:"dir,omitempty"  yaml:"dir,omitempty" toml:"dir,omitempty"`
	// 日志文件路径（模式），可以用多级路径，拼接在 Dir 后面。
	// 可用时间格式、事件格式，如 %Y-%m-%d/${level}.log 表示 按日期和级别分割
	FilePattern string `mapstructure:"filePattern,omitempty" json:"filePattern,omitempty" yaml:"filePattern,omitempty" toml:"filePattern,omitempty"`
	// 日志留存时间（天）
	MaxAge int `mapstructure:"maxAge,omitempty" json:"maxAge,omitempty" yaml:"maxAge,omitempty" toml:"maxAge,omitempty"`

	// 输出格式 console(默认) 或者 json
	Format string `mapstructure:"format,omitempty" json:"format,omitempty" yaml:"format,omitempty" toml:"format,omitempty"`

	// 日志前缀
	//Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`

	// 编码级，可用值为:
	//  color
	//  capital
	//  capitalColor
	EncodeLevel string `mapstructure:"encodeLevel,omitempty" json:"encodeLevel,omitempty" yaml:"encodeLevel,omitempty" toml:"encodeLevel,omitempty"`

	// 自定义特殊字段的 key，主要是用在 json 格式中
	LogKey *KeyLogConfig `mapstructure:"logKey,omitempty" json:"logKey,omitempty" yaml:"logKey,omitempty" toml:"logKey,omitempty"`

	// 显示行号
	ShowLine bool `mapstructure:"showLine,omitempty" json:"showLine,omitempty" yaml:"showLine,omitempty" toml:"showLine,omitempty"`
	// 输出到文件同时输出到控制台
	AppendConsole bool `mapstructure:"appendConsole,omitempty" json:"appendConsole,omitempty" yaml:"appendConsole,omitempty" toml:"appendConsole,omitempty"`
}

// FixDefault 设置默认值
func (c *LogConfig) FixDefault() {
	if c.Level == "" {
		c.Level = "info"
	}
	c.EncodeLevel = strings.TrimSpace(c.EncodeLevel)
	c.Dir = strings.TrimSpace(c.Dir)
	c.FilePattern = strings.TrimSpace(c.FilePattern)
	if c.Dir != "" || c.FilePattern != "" {
		if c.FilePattern == "" {
			c.FilePattern = DefaultFilePattern
		}
	}
}

// IsFileEnabled 是否启用文件输出
func (c *LogConfig) IsFileEnabled() bool {
	return c.Dir != "" || c.FilePattern != ""
}

// ISFileSplitByEvent 是否按事件分割文件
func (c *LogConfig) ISFileSplitByEvent() bool {
	return c.IsFileEnabled() && strings.Contains(c.FilePattern, HolderLevel)
}
