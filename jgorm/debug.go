package jgorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"path"
	"strings"
	_ "unsafe"
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

// SafeStr 限定字符串长度，避免字段溢出
func SafeStr(s string, size int) string {
	if size <= 1 || len(s) < size {
		return s
	}
	cc := []rune(s)
	if size <= 1 || len(cc) < size {
		return s
	}
	return string(cc[:size])
}

// SetDbLogger 打开调试日志
func SetDbLogger(db *gorm.DB, logSql string) *gorm.DB {
	if logSql == "1" || logSql == "true" || logSql == "info" {
		db = db.Debug()
	} else {
		var t logger.LogLevel
		switch strings.ToLower(logSql) {
		case "debug", "info":
			t = logger.Info
		case "warn", "warning":
			t = logger.Warn
		case "error":
			t = logger.Error
		}
		if t > 0 {
			db = db.Session(&gorm.Session{Logger: db.Logger.LogMode(t)})
		}
	}
	return db
}
