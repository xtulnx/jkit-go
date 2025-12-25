package driver_sqlite

import (
	"github.com/xtulnx/jkit-go/jgorm/builder"
	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite" // 纯 golang 实现
	"gorm.io/gorm"
)

func init() {
	builder.DefaultGormBuilder.
		Register("sqlite", func(dsn string) gorm.Dialector {
			return sqlite.Open(dsn)
		})
}
