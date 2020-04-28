package getzap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

// GetDevelopmentLogger returns a sophisticated, sophisticated logger for development.
//
// Logs with level below than ERROR will be written to:
//    a) stdout in TSV format with colored level;
//    b) normal log files in JSON format with rotation and expiration;
//
// Logs with level ERROR or above will be written to:
//    a) stderr in TSV format with colored level;
//    b) error log files in JSON format with rotation and expiration;
//
// Logs with level DPANIC or above will cause panic after writing the message.
func GetDevelopmentLogger(normalLogPath, errorLogPath string) *zap.Logger {
	var (
		consoleEncoder = zapcore.NewConsoleEncoder(*genEncoderConfig(zapcore.CapitalColorLevelEncoder))
		jsonEncoder    = zapcore.NewJSONEncoder(*genEncoderConfig(zapcore.LowercaseLevelEncoder))

		writeStdout = zapcore.AddSync(os.Stdout)
		writeStderr = zapcore.AddSync(os.Stderr)

		normalLogFile = zapcore.AddSync(&lj.Logger{
			Filename:   normalLogPath,
			MaxBackups: 20, // numbers
			MaxSize:    10, // megabytes
			MaxAge:     14, // days
		})
		errorLogFile = zapcore.AddSync(&lj.Logger{
			Filename:   errorLogPath,
			MaxBackups: 20, // numbers
			MaxSize:    10, // megabytes
			MaxAge:     14, // days
		})
	)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, writeStdout, normalLevel),
		zapcore.NewCore(jsonEncoder, normalLogFile, normalLevel),
		zapcore.NewCore(consoleEncoder, writeStderr, errorLevel),
		zapcore.NewCore(jsonEncoder, errorLogFile, normalLevel),
	)

	options := []zap.Option{
		zap.Development(),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}

	return zap.New(core, options...)
}

// GetProductionLogger returns a sophisticated, customized logger for production deployment.
//
// Logs with level below than INFO will be ignored.
//
// Logs with level INFO or above will be written to:
//    a) stdout in JSON format;
//    b) log files in JSON format with rotation and compression;
//
// Logs with level PANIC or above will cause panic after writing the message.
func GetProductionLogger() *zap.Logger {
	return nil
}

var (
	genEncoderConfig = func(lvlEnc zapcore.LevelEncoder) *zapcore.EncoderConfig {
		return &zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    lvlEnc,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}
	// check if the log level is less than error
	normalLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	// check if the level is greater than or equal to error
	errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// check if the level is greater than or equal to info
	infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
)
