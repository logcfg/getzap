package getzap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

func GetDevColorConsoleLogger() *zap.Logger {
	return getConsoleLogger(
		zapcore.NewConsoleEncoder(
			*customEncoderConfig(zapcore.CapitalColorLevelEncoder,
				zapcore.ISO8601TimeEncoder)),
		true)
}

func GetDevJsonConsoleLogger() *zap.Logger {
	return getConsoleLogger(
		zapcore.NewJSONEncoder(
			*customEncoderConfig(zapcore.LowercaseLevelEncoder,
				zapcore.ISO8601TimeEncoder)),
		true)
}

func GetProdJsonConsoleLogger() *zap.Logger {
	return getConsoleLogger(
		zapcore.NewJSONEncoder(
			*customEncoderConfig(zapcore.LowercaseLevelEncoder,
				zapcore.ISO8601TimeEncoder)),
		false)
}

func GetProdJsonConsoleAndFileLogger(logPath string) *zap.Logger {
	return getConsoleAndFileLogger(logPath)
}

var (
	customEncoderConfig = func(lvlEnc zapcore.LevelEncoder, timeEnc zapcore.TimeEncoder) *zapcore.EncoderConfig {
		return &zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    lvlEnc,
			EncodeTime:     timeEnc,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}
	// check if the level is greater than or equal to info
	geInfoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	// check if the level is greater than or equal to error
	geErrorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// check if the level is less than error
	ltErrorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	// for all level
	allLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})
)

func getConsoleLogger(encoder zapcore.Encoder, isDev bool) *zap.Logger {
	var (
		writeStdout = zapcore.AddSync(os.Stdout)
		writeStderr = zapcore.AddSync(os.Stderr)
	)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writeStdout, ltErrorLevel),
		zapcore.NewCore(encoder, writeStderr, geErrorLevel),
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

func getConsoleAndFileLogger(logPath string) *zap.Logger {
	var (
		jsonTimefmtEncoder = zapcore.NewJSONEncoder(*customEncoderConfig(zapcore.LowercaseLevelEncoder, zapcore.ISO8601TimeEncoder))
		jsonEpochEncoder   = zapcore.NewJSONEncoder(*customEncoderConfig(zapcore.LowercaseLevelEncoder, zapcore.EpochTimeEncoder))

		writeStdout = zapcore.AddSync(os.Stdout)
		writeStderr = zapcore.AddSync(os.Stderr)

		logFile = zapcore.AddSync(&lj.Logger{
			Filename:   logPath,
			MaxBackups: 200, // numbers
			MaxSize:    5,   // megabytes
			MaxAge:     30,  // days
			Compress:   true,
		})
	)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEpochEncoder, logFile, geInfoLevel),
		zapcore.NewCore(jsonTimefmtEncoder, writeStdout, ltErrorLevel),
		zapcore.NewCore(jsonTimefmtEncoder, writeStderr, geErrorLevel),
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}

	return zap.New(core, options...)
}
