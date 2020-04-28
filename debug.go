package getzap

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	plainEpochTimeEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	colorTimeFormatEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	normalLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
)

func GetDevColorConsoleLogger() *zap.Logger {
	return getConsoleLogger(zapcore.NewConsoleEncoder(colorTimeFormatEncoderConfig), true)
}

func GetDevJsonConsoleLogger() *zap.Logger {
	return getConsoleLogger(zapcore.NewJSONEncoder(plainEpochTimeEncoderConfig), true)
}

func GetProdJsonConsoleLogger() *zap.Logger {
	return getConsoleLogger(zapcore.NewJSONEncoder(plainEpochTimeEncoderConfig), false)
}

func GetProdJsonConsoleAndFileLogger() *zap.Logger {
	// TODO:
	return getConsoleLogger(zapcore.NewJSONEncoder(plainEpochTimeEncoderConfig), true)
}

func getConsoleLogger(encoder zapcore.Encoder, isDev bool) *zap.Logger {
	var (
		writeStdout = zapcore.AddSync(os.Stdout)
		writeStderr = zapcore.AddSync(os.Stderr)
	)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writeStdout, normalLevel),
		zapcore.NewCore(encoder, writeStderr, errorLevel),
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	if isDev {
		options = append(options, zap.Development())
	}

	return zap.New(core, options...)
}

func learn_global() {
	jsonEncoderConfig := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	consoleEncoderConfig := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	errorFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "log/error.log",
		MaxBackups: 0, // numbers
		MaxSize:    5, // megabytes
		MaxAge:     7, // days
		Compress:   true,
	})
	normalFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "log/access.log",
		MaxBackups: 0, // numbers
		MaxSize:    5, // megabytes
		MaxAge:     7, // days
		Compress:   true,
	})

	writeStderr := zapcore.AddSync(os.Stderr)
	writeStdout := zapcore.AddSync(os.Stdout)

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	normalLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	host, _ := os.Hostname()
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoderConfig, normalFile, normalLevel),
		zapcore.NewCore(jsonEncoderConfig, errorFile, errorLevel),
		zapcore.NewCore(consoleEncoderConfig, writeStdout, normalLevel),
		zapcore.NewCore(consoleEncoderConfig, writeStderr, errorLevel),
	)
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.DPanicLevel)).With(zap.String("host", host))

	sugar := logger.Sugar()
	defer sugar.Sync()

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2), zap.Duration("sleep", 233*time.Millisecond))
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")
	logger.DPanic("This is a DPANIC message")

	_ = zap.ReplaceGlobals(logger)
}
