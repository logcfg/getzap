package getzap

import (
	"os"
	"strings"

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
//    b) error log files in JSON format with rotation (max size is 10MB, retain at most 20 files) and expiration (14 days);
//
// Logs with level DPANIC or above will cause panic after writing the message.
//
// Log files won't be created if the given log file path is empty or blank.
func GetDevelopmentLogger(normalLogPath, errorLogPath string) *zap.Logger {
	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, writeStdout, normalLevel),
		zapcore.NewCore(consoleEncoder, writeStderr, errorLevel),
	}

	normalLogPath, errorLogPath = strings.TrimSpace(normalLogPath), strings.TrimSpace(errorLogPath)
	if len(normalLogPath) > 0 {
		normalLogFile := zapcore.AddSync(&lj.Logger{
			Filename:   normalLogPath,
			MaxBackups: 20, // numbers
			MaxSize:    10, // megabytes
			MaxAge:     14, // days
		})
		cores = append(cores, zapcore.NewCore(jsonEncoder, normalLogFile, normalLevel))
	}
	if len(errorLogPath) > 0 {
		errorLogFile := zapcore.AddSync(&lj.Logger{
			Filename:   errorLogPath,
			MaxBackups: 20, // numbers
			MaxSize:    10, // megabytes
			MaxAge:     14, // days
		})
		cores = append(cores, zapcore.NewCore(jsonEncoder, errorLogFile, errorLevel))
	}

	options := []zap.Option{
		zap.Development(),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}

	return zap.New(zapcore.NewTee(cores...), options...)
}

// GetProductionLogger returns a sophisticated, customized logger for production deployment.
//
// Logs with level below than INFO will be ignored.
//
// Logs with level INFO or above will be written to:
//    a) stdout in JSON format;
//    b) log files in JSON format with rotation (max size is 10MB), expiration (30 days) and compression (.gz);
//
// Logs with level PANIC or above will cause panic after writing the message.
//
// Log files won't be created if the given log file path is empty or blank.
func GetProductionLogger(logPath string) *zap.Logger {
	cores := []zapcore.Core{
		zapcore.NewCore(jsonEncoder, writeStdout, infoLevel),
	}

	if logPath = strings.TrimSpace(logPath); len(logPath) > 0 {
		logFile := zapcore.AddSync(&lj.Logger{
			Filename: logPath,
			MaxSize:  10, // megabytes
			MaxAge:   30, // days
			Compress: true,
		})
		cores = append(cores, zapcore.NewCore(jsonEncoder, logFile, infoLevel))
	}

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}

	return zap.New(zapcore.NewTee(cores...), options...)
}

var (
	getEncoderConfig = func(lvlEnc zapcore.LevelEncoder) *zapcore.EncoderConfig {
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
	// log encoders
	consoleEncoder = zapcore.NewConsoleEncoder(*getEncoderConfig(zapcore.CapitalColorLevelEncoder))
	jsonEncoder    = zapcore.NewJSONEncoder(*getEncoderConfig(zapcore.LowercaseLevelEncoder))
	// std log writers
	writeStdout = zapcore.AddSync(os.Stdout)
	writeStderr = zapcore.AddSync(os.Stderr)
)
