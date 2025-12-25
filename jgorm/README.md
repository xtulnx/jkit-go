## gorm 辅助

### 配置 config

数据源规则参考：  https://gorm.io/zh_CN/docs/connecting_to_the_database.html

示例:

```go
type ServerConfig struct {
    MySqlDb GeneralDB `mapstructure:"mysql" json:"mysql" yaml:"mysql"` // 指定是 mysql
    Db SpecializedDB  `mapstructure:"db" json:"db" yaml:"db"`          // 通用配置，需要在配置中指定 type
}
```

配置内容

```yaml
mysql:
  path: 127.0.0.1
  port: "3306"
  config: charset=utf8mb4&parseTime=True&loc=Asia%2fShanghai
  db-name: cms-dev
  username: cms-user
  password: 0000000000000000
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: "info"
  log-zap: false
db:
  type: sqlite
  dsn: file::memory:?cache=shared
  log-mode: "debug"
  log-zap: "warn"
```

### 初始化 builder

需要 先引用 driver_default 加载驱动


### 日志

* log-mode:

* log-zap: 是否输出到 zap 日志管理器，需要先调用 debug_zap

    .SetDefaultZapLogger(l *zap.Logger, level zapcore.Level) 指定 zap.logger


