package builder

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xtulnx/jkit-go/jgorm/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DialectorFunc 是 GORM dialector 的工厂函数
type DialectorFunc func(dsn string) gorm.Dialector

// LoggerBuilder 创建日志
type LoggerBuilder func(c config.GeneralDB) logger.Interface

// GormBuilder 构造器
type GormBuilder struct {
	drivers map[string]DialectorFunc
}

// NewGormBuilder 创建空构造器
func NewGormBuilder() *GormBuilder {
	return &GormBuilder{
		drivers: make(map[string]DialectorFunc),
	}
}

// Register 注册数据库方言支持
func (b *GormBuilder) Register(dialect string, fn DialectorFunc) *GormBuilder {
	b.drivers[dialect] = fn
	return b
}

// ToLogLevel 将字符串日志模式转为 GORM LogLevel
func ToLogLevel(logMode string) logger.LogLevel {
	switch strings.ToLower(logMode) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info", "debug":
		return logger.Info
	case "1", "true":
		return logger.Info
	case "0", "false":
		return logger.Error
	default:
		return logger.Info
	}
}

// Build 根据 DsnProvider 创建 *gorm.DB
func (b *GormBuilder) Build(provider config.DsnProvider) (*gorm.DB, error) {
	dsn := provider.Dsn()
	if dsn == "" {
		return nil, ErrEmptyDSN
	}

	dialect := provider.Dialect()
	fn, ok := b.drivers[dialect]
	if !ok {
		return nil, ErrUnsupportedDialect{Dialect: dialect}
	}

	general, hasGeneral := config.ExtractGeneralDB(provider)

	dialector := fn(dsn)
	if dialector == nil {
		return nil, ErrUnsupportedDialect{Dialect: dialect, Msg: "invalid DSN"}
	}

	// 日志配置
	var gormLogger logger.Interface
	if defaultLoggerBuilder != nil {
		gormLogger = defaultLoggerBuilder(general)
	} else {
		if hasGeneral {
			gormLogger = logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: 200 * time.Millisecond,
					LogLevel:      ToLogLevel(general.LogMode),
					Colorful:      true,
				},
			)
		} else {
			gormLogger = logger.Default
		}
	}

	// 命名策略
	var namingStrategy schema.Namer
	if hasGeneral {
		if general.Prefix != "" || general.Singular {
			namingStrategy = &schema.NamingStrategy{
				TablePrefix:         general.Prefix,
				SingularTable:       general.Singular,
				NameReplacer:        nil,
				NoLowerCase:         false,
				IdentifierMaxLength: 0,
			}
		}
	}

	gormConfig := &gorm.Config{
		Logger:         gormLogger,
		NamingStrategy: namingStrategy,
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, err
	}

	// 应用通用配置（连接池等）
	if hasGeneral {
		switch dialect {
		case "mysql":
			if general.Engine != "" {
				db.InstanceSet("gorm:table_options", "ENGINE="+general.Engine)
			} else {
				db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
			}
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		if general.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(general.MaxIdleConns)
		}
		if general.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(general.MaxOpenConns)
		}
	}
	return db, nil
}

// 预定义构造器

var (
	// DefaultGormBuilder 支持所有主流数据库（按需导入驱动）
	DefaultGormBuilder = NewGormBuilder()

	// LiteGormBuilder 仅支持 SQLite（最小依赖）
	LiteGormBuilder = NewGormBuilder()

	// 默认的日志输出
	defaultLoggerBuilder LoggerBuilder = nil
)

// SetDefaultWriter 指定默认日志输出
func SetDefaultLoggerBuilder(b LoggerBuilder) {
	defaultLoggerBuilder = b
}

// 注意：预定义构造器的实际注册需在 init() 或由用户显式调用，
// 因为 Go 不允许在全局变量中直接引用可能未导入的包。
// 所以我们提供辅助函数来“激活”它们。

// 自定义错误
var ErrEmptyDSN = errors.New("DSN is empty")

type ErrUnsupportedDialect struct {
	Dialect string
	Msg     string
}

func (e ErrUnsupportedDialect) Error() string {
	s := "unsupported dialect: " + e.Dialect
	if e.Msg == "" {
		return s
	}
	return s + "," + e.Msg
}
