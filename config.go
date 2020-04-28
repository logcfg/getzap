package getzap

import (
	"go.uber.org/zap"
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
func GetDevelopmentLogger() *zap.Logger {
	return nil
}

// GetProductionLogger returns a sophisticated, customized logger for production deployment.
//
// Logs with level INFO or above will be written to:
//    a) stdout in JSON format;
//    b) log files in JSON format with rotation and compression;
//
// Logs with level PANIC or above will cause panic after writing the message.
func GetProductionLogger() *zap.Logger {
	return nil
}
