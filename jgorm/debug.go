package jgorm

import (
	"path"
	"strings"
	_ "unsafe"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

//go:linkname gormSourceDir gorm.io/gorm/utils.gormSourceDir
var gormSourceDir string

func DebugGorm() {
	_ = utils.FileWithLineNum()
	//fmt.Println("sourceDir", gormSourceDir)
	if gormSourceDir != "" {
		gormSourceDir = path.Dir(path.Dir(gormSourceDir))
		//fmt.Println("targetDir", gormSourceDir)
	}
}

// LoggerLevel 日志等级，如 info、warn、error、silent。
func LoggerLevel(logSql string) logger.LogLevel {
	var t logger.LogLevel
	switch strings.ToLower(logSql) {
	case "1", "true":
		t = logger.Info
	case "debug", "info":
		t = logger.Info
	case "warn", "warning":
		t = logger.Warn
	case "error":
		t = logger.Error
	case "silent":
		t = logger.Silent
	}
	return t
}

// SetDbLogger 打开调试日志。
// logSql 是日志等级，如 info、warn、error、silent。
func SetDbLogger(db *gorm.DB, logSql string) *gorm.DB {
	l := LoggerLevel(logSql)
	if l > 0 {
		db = db.Session(&gorm.Session{Logger: db.Logger.LogMode(l)})
	}
	return db
}

// ShowSQL 获取 db1 的 sql
func ShowSQL(db1 *gorm.DB) string {
	return db1.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var tm interface{}
		return tx.Scan(&tm)
	})
}
