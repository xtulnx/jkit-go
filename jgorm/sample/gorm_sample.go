package sample

import (
	"github.com/xtulnx/jkit-go/jgorm/builder"
	"github.com/xtulnx/jkit-go/jgorm/config"
	"github.com/xtulnx/jkit-go/jgorm/debug_zap"
	_ "github.com/xtulnx/jkit-go/jgorm/driver_default"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

func InitEnv(logger *zap.Logger) {
	debug_zap.SetDefaultZapLogger(logger, zapcore.InfoLevel)
}

func OpenDb(c config.DsnProvider) (*gorm.DB, error) {
	return builder.DefaultGormBuilder.Build(c)
}
