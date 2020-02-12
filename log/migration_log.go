package log

import (
	"fmt"
	"sync"

	"github.com/TRON-US/btfs-migration-toolkit/core"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once sync.Once
)

// Init zap log.
func initLog() {
	logger = InitLogger(fmt.Sprintf("%s/%s", core.Conf.Logger.Path, "migration.log"), core.Conf.Logger.Level)
}

// Get zap log instance.
func Logger() *zap.Logger {
	once.Do(func() {
		initLog()
	})
	return logger
}
