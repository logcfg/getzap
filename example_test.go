package getzap_test

import (
	"time"

	"github.com/logcfg/getzap"
	"go.uber.org/zap"
)

func ExampleGetDevelopmentLogger() {
	logger := getzap.GetDevelopmentLogger("log/access.log", "log/error.log")
	defer logger.Sync()

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.Int("id", 1), zap.Duration("sleep", 233*time.Millisecond))
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")

	// Output:
}

func ExampleGetProductionLogger() {
	logger := getzap.GetProductionLogger("log/app.log")
	defer logger.Sync()

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.Int("id", 1), zap.Duration("sleep", 233*time.Millisecond))
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")
	logger.DPanic("This is a DPANIC message")

	// Output:
}
