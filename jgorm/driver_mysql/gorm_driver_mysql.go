package driver_mysql

import (
	"github.com/xtulnx/jkit-go/jgorm/builder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	builder.DefaultGormBuilder.
		Register("mysql", func(dsn string) gorm.Dialector {
			return mysql.Open(dsn) // 需要 import _ "gorm.io/driver/mysql"
		})
}
