package log

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once sync.Once
)

// Init zap log.
func initLog() {
	logPath := fmt.Sprintf("log.path")
	logLevel := fmt.Sprintf("log.level")
	logger = InitLogger(fmt.Sprintf("%s/%s", viper.GetString(logPath), "migration.log"), viper.GetString(logLevel))
}

// Get zap log instance.
func Logger() *zap.Logger {
	once.Do(func() {
		initLog()
	})
	return logger
}
