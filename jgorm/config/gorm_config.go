package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strings"
)

// DsnProvider 接口：任何可提供 DSN 和方言的类型需实现此接口
type DsnProvider interface {
	Dsn() string
	Dialect() string
}

// GeneralDB 是所有数据库类型的公共配置字段
type GeneralDB struct {
	FullDsn string `mapstructure:"dsn" json:"dsn" yaml:"dsn"` // 完整 DSN，若提供则优先使用

	Path     string `mapstructure:"path" json:"path" yaml:"path"`             // 地址（主机或文件路径）
	Port     string `mapstructure:"port" json:"port" yaml:"port"`             // 端口
	Config   string `mapstructure:"config" json:"config" yaml:"config"`       // 额外连接参数
	Dbname   string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`    // 数据库名
	Username string `mapstructure:"username" json:"username" yaml:"username"` // 用户名
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码

	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 表前缀
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`                   // 禁用复数表名
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"`        // 引擎（仅 MySQL 有效）
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 最大空闲连接
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 最大打开连接

	LogMode string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"` // GORM 日志级别
	LogZap  string `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`    // 是否用 zap 记录日志，默认 true 用 info，如 warn、debug、error 等
}

// DsnBuilder 定义 DSN 构建策略
type DsnBuilder interface {
	Build(g GeneralDB) string
}

// --- 各数据库的 DSN 构建器 ---
var (
	mysqlBuilderInstance  = mysqlBuilder{}
	mssqlBuilderInstance  = mssqlBuilder{}
	pgsqlBuilderInstance  = pgsqlBuilder{}
	oracleBuilderInstance = oracleBuilder{}
	sqliteBuilderInstance = sqliteBuilder{}
)

var dsnBuilders = map[string]DsnBuilder{
	"mysql":     mysqlBuilderInstance,
	"mssql":     mssqlBuilderInstance,
	"sqlserver": mssqlBuilderInstance,
	"pgsql":     pgsqlBuilderInstance,
	"postgres":  pgsqlBuilderInstance,
	"oracle":    oracleBuilderInstance,
	"sqlite":    sqliteBuilderInstance,
}

// RegDsnBuilder 注册自定义的 DSN 生成器，可以覆盖已有的
func RegDsnBuilder(typ string, builder DsnBuilder) {
	dsnBuilders[typ] = builder
}

type mysqlBuilder struct{}

func (mysqlBuilder) Build(g GeneralDB) string {
	return g.Username + ":" + g.Password + "@tcp(" + g.Path + ":" + g.Port + ")/" + g.Dbname + "?" + g.Config
}

type mssqlBuilder struct{}

func (mssqlBuilder) Build(g GeneralDB) string {
	return "sqlserver://" + g.Username + ":" + g.Password + "@" + g.Path + ":" + g.Port + "?database=" + g.Dbname + "&encrypt=disable"
}

type pgsqlBuilder struct{}

func (pgsqlBuilder) Build(g GeneralDB) string {
	return "host=" + g.Path + " user=" + g.Username + " password=" + g.Password + " dbname=" + g.Dbname + " port=" + g.Port + " " + g.Config
}

type oracleBuilder struct{}

func (oracleBuilder) Build(g GeneralDB) string {
	return fmt.Sprintf("oracle://%s:%s@%s/%s?%s",
		url.PathEscape(g.Username),
		url.PathEscape(g.Password),
		net.JoinHostPort(g.Path, g.Port),
		url.PathEscape(g.Dbname),
		g.Config,
	)
}

type sqliteBuilder struct{}

func (sqliteBuilder) Build(g GeneralDB) string {
	return filepath.Join(g.Path, g.Dbname+".db")
}

// --- 具体数据库类型（保留用于强类型配置） ---

type Mysql struct{ GeneralDB }

func (m *Mysql) Dsn() string {
	if m.FullDsn != "" {
		return m.FullDsn
	}
	return mysqlBuilderInstance.Build(m.GeneralDB)
}
func (m *Mysql) Dialect() string { return "mysql" }

type Mssql struct{ GeneralDB }

func (m *Mssql) Dsn() string {
	if m.FullDsn != "" {
		return m.FullDsn
	}
	return mssqlBuilderInstance.Build(m.GeneralDB)
}
func (m *Mssql) Dialect() string { return "sqlserver" } // GORM 使用 sqlserver

type Pgsql struct{ GeneralDB }

func (p *Pgsql) Dsn() string {
	if p.FullDsn != "" {
		return p.FullDsn
	}
	return pgsqlBuilderInstance.Build(p.GeneralDB)
}
func (p *Pgsql) Dialect() string { return "postgres" } // GORM 使用 postgres

// LinkDsn 保留 Pgsql 特有方法（如需创建其他 DB）
func (p *Pgsql) LinkDsn(dbname string) string {
	if p.FullDsn != "" {
		return p.FullDsn
	}
	g := p.GeneralDB
	g.Dbname = dbname
	return pgsqlBuilderInstance.Build(g)
}

type Oracle struct{ GeneralDB }

func (o *Oracle) Dsn() string {
	if o.FullDsn != "" {
		return o.FullDsn
	}
	return oracleBuilderInstance.Build(o.GeneralDB)
}
func (o *Oracle) Dialect() string { return "oracle" }

type Sqlite struct{ GeneralDB }

func (s *Sqlite) Dsn() string {
	if s.FullDsn != "" {
		return s.FullDsn
	}
	return sqliteBuilderInstance.Build(s.GeneralDB)
}
func (s *Sqlite) Dialect() string { return "sqlite" }

// --- SpecializedDB：用于动态/通用数据源配置 ---

type SpecializedDB struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`

	Type      string `mapstructure:"type" json:"type" yaml:"type"` // mysql, sqlite, mssql, pgsql, oracle
	AliasName string `mapstructure:"alias-name" json:"alias-name" yaml:"alias-name"`
	Disable   bool   `mapstructure:"disable" json:"disable" yaml:"disable"`
}

func (s *SpecializedDB) Dsn() string {
	if s.FullDsn != "" {
		return s.FullDsn
	}
	if builder, ok := dsnBuilders[strings.ToLower(s.Type)]; ok {
		return builder.Build(s.GeneralDB)
	}
	return ""
}

// Dialect 返回 GORM 兼容的驱动名
func (s *SpecializedDB) Dialect() string {
	switch strings.ToLower(s.Type) {
	case "pgsql":
		return "postgres"
	case "mssql":
		return "sqlserver"
	default:
		return strings.ToLower(s.Type)
	}
}

// Validate 验证配置是否合法
func (s *SpecializedDB) Validate() error {
	if s.Disable {
		return nil
	}
	typ := strings.ToLower(s.Type)
	if typ == "" {
		return errors.New("database type is required")
	}
	if _, ok := dsnBuilders[typ]; !ok {
		return fmt.Errorf("unsupported database type: %s", s.Type)
	}
	if s.FullDsn == "" {
		if typ == "sqlite" {
			if s.Path == "" {
				return errors.New("sqlite requires 'path'")
			}
		} else {
			if s.Username == "" || s.Dbname == "" || s.Path == "" {
				return fmt.Errorf("missing required fields (username, dbname, path) for %s", s.Type)
			}
			if s.Port == "" {
				return fmt.Errorf("port is required for %s", s.Type)
			}
		}
	}
	return nil
}

type DbProvider interface {
	GetGeneralDB() GeneralDB
}

// ExtractGeneralDB 辅助：从 DsnProvider 中提取 GeneralDB（用于配置复用）
func ExtractGeneralDB(p DsnProvider) (GeneralDB, bool) {
	switch v := p.(type) {
	case *Mysql:
		return v.GeneralDB, true
	case *Pgsql:
		return v.GeneralDB, true
	case *Mssql:
		return v.GeneralDB, true
	case *Oracle:
		return v.GeneralDB, true
	case *Sqlite:
		return v.GeneralDB, true
	case *SpecializedDB:
		return v.GeneralDB, true
	case DbProvider:
		return v.GetGeneralDB(), true
	default:
		return GeneralDB{}, false
	}
}

func NewDbProvider1(dialect, dsn string) DsnProvider {
	return &SpecializedDB{Type: dialect, GeneralDB: GeneralDB{FullDsn: dsn}}
}
