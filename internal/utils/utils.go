package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the level of logging
type LogLevel int

const (
	// LogLevelDebug logs everything
	LogLevelDebug LogLevel = iota
	// LogLevelInfo logs info, warn, error
	LogLevelInfo
	// LogLevelWarn logs warn, error
	LogLevelWarn
	// LogLevelError logs only errors
	LogLevelError
	// LogLevelNone disables logging
	LogLevelNone
)

var (
	// Logger is the global logger instance
	Logger *CustomLogger
	// CurrentLogLevel controls the verbosity of logging
	CurrentLogLevel = LogLevelInfo
)

// CustomLogger wraps the standard logger with log level functionality
type CustomLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

// InitLogger creates a new logger with the specified log level
func InitLogger(level LogLevel) {
	CurrentLogLevel = level

	debugHandle := os.Stdout
	infoHandle := os.Stdout
	warnHandle := os.Stdout
	errorHandle := os.Stderr

	Logger = &CustomLogger{
		debugLogger: log.New(debugHandle, "[DEBUG] ", log.Ldate|log.Ltime),
		infoLogger:  log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime),
		warnLogger:  log.New(warnHandle, "[WARN] ", log.Ldate|log.Ltime),
		errorLogger: log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime),
	}
}

// Debug logs debug messages
func (l *CustomLogger) Debug(format string, v ...interface{}) {
	if CurrentLogLevel <= LogLevelDebug {
		msg := fmt.Sprintf(format, v...)
		l.debugLogger.Println(msg)
	}
}

// Info logs info messages
func (l *CustomLogger) Info(format string, v ...interface{}) {
	if CurrentLogLevel <= LogLevelInfo {
		msg := fmt.Sprintf(format, v...)
		l.infoLogger.Println(msg)
	}
}

// Warn logs warning messages
func (l *CustomLogger) Warn(format string, v ...interface{}) {
	if CurrentLogLevel <= LogLevelWarn {
		msg := fmt.Sprintf(format, v...)
		l.warnLogger.Println(msg)
	}
}

// Error logs error messages
func (l *CustomLogger) Error(format string, v ...interface{}) {
	if CurrentLogLevel <= LogLevelError {
		msg := fmt.Sprintf(format, v...)
		l.errorLogger.Println(msg)
	}
}

// TimeTrack logs the execution time of a function
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	if Logger != nil {
		Logger.Debug("%s took %s", name, elapsed)
	} else {
		fmt.Printf("%s took %s\n", name, elapsed)
	}
}
