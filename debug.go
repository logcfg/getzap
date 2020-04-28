package getzap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	jsonEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
	consoleEncoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
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
	aboveInfoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	aboveErrorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	belowErrorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
)

func GetDevColorConsoleLogger() *zap.Logger {
	return getConsoleLogger(consoleEncoder, true)
}

func GetDevJsonConsoleLogger() *zap.Logger {
	return getConsoleLogger(jsonEncoder, true)
}

func GetProdJsonConsoleLogger() *zap.Logger {
	return getConsoleLogger(jsonEncoder, false)
}

func GetProdJsonConsoleAndFileLogger(logPath string) *zap.Logger {
	return getConsoleAndFileLogger(jsonEncoder, logPath)
}

func getConsoleLogger(encoder zapcore.Encoder, isDev bool) *zap.Logger {
	var (
		writeStdout = zapcore.AddSync(os.Stdout)
		writeStderr = zapcore.AddSync(os.Stderr)
	)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writeStdout, belowErrorLevel),
		zapcore.NewCore(encoder, writeStderr, aboveErrorLevel),
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

func getConsoleAndFileLogger(encoder zapcore.Encoder, logPath string) *zap.Logger {
	var (
		writeStdout = zapcore.AddSync(os.Stdout)
		writeStderr = zapcore.AddSync(os.Stderr)
		logFile     = zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxBackups: 200, // numbers
			MaxSize:    5,   // megabytes
			MaxAge:     30,  // days
			Compress:   true,
		})
	)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, logFile, aboveInfoLevel),
		zapcore.NewCore(encoder, writeStdout, belowErrorLevel),
		zapcore.NewCore(encoder, writeStderr, aboveErrorLevel),
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	return zap.New(core, options...)
}
