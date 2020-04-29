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
	logger.Info("This is an INFO message with fields", zap.Int("id", 1), zap.Duration("sleep", 128*time.Millisecond))
	logger.Sugar().Infof("The answer to life the universe and everything = %d", 42)
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")

	// Output:
}

func ExampleGetProductionLogger() {
	logger := getzap.GetProductionLogger("log/app.log")
	defer logger.Sync()

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.Int("id", 1), zap.Duration("sleep", 64*time.Millisecond))
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")
	logger.DPanic("This is a DPANIC message")

	// Output:
}

func ExampleSetGlobalDevelopmentLogger() {
	getzap.SetGlobalDevelopmentLogger("", "log/error.log")
	logger := zap.L().Sugar().Named("dev")

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Infow("This is an INFO message with fields", "id", 1, "sleep", 256*time.Millisecond)
	logger.Infof("The answer to life the universe and everything = %d", 42)
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")

	// Output:
}

func ExampleSetGlobalProductionLogger() {
	getzap.SetGlobalProductionLogger("log/app.log")
	logger := zap.L().Named("prod")

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.Int("id", 1), zap.Duration("sleep", 32*time.Millisecond))
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")
	logger.DPanic("This is a DPANIC message")

	// Output:
}
