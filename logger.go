// Package logger provides a structured logging solution.
// It supports multiple output modes (terminal, file, both) and various log levels.
package gologger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Output modes for logger configuration.
const (
	OutputTerminal = "terminal"
	OutputFile     = "file"
	OutputBoth     = "both"
)

// Log levels for logger configuration.
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// Context key for request ID.
type contextKey string

const (
	RequestIDKey contextKey = "gologger-request-id"
)

// Logger provides a simplified structured logging interface.
type Logger struct {
	log          *zap.SugaredLogger
	ctx          context.Context
	level        string
	message      string
	data         []any
	hasData      bool
	requestIDKey string // Custom key for request ID in logs
}

// LoggerConfig holds configuration options for the logger.
type LoggerConfig struct {
	OutputMode   string // Output mode: OutputTerminal, OutputFile, or OutputBoth
	LogLevel     string // Log level: LevelDebug, LevelInfo, LevelWarn, or LevelError
	LogDir       string // Directory for log files
	RequestIDKey string // Custom key for request ID in logs (default: "request-id")
}

// NewLogger creates a new Logger instance with default configuration.
// Default settings: output to both terminal and file, debug level, logs saved to "logger" directory.
func NewLogger() Logger {
	return NewLoggerWithConfig(LoggerConfig{
		OutputMode:   OutputBoth,   // default: both terminal and file
		LogLevel:     LevelDebug,   // default: debug level
		LogDir:       "logger",     // default: logger directory
		RequestIDKey: "request-id", // default: request-id key
	})
}

// NewLoggerWithConfig creates a new Logger instance with custom configuration.
func NewLoggerWithConfig(config LoggerConfig) Logger {
	// Set default request ID key if not provided
	requestIDKey := config.RequestIDKey
	if requestIDKey == "" {
		requestIDKey = "request-id"
	}

	return Logger{
		log:          initLogWithConfig(config),
		ctx:          context.Background(),
		level:        "",
		message:      "",
		data:         make([]any, 0),
		hasData:      false,
		requestIDKey: requestIDKey,
	}
}

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID retrieves the request ID from the context.
// Returns empty string if no request ID is found.
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// prefix generates a log file prefix with current date.
func prefix() string {
	return "logger-" + time.Now().Format("2006-01-02")
}

// initLogWithConfig creates a logger with custom configuration.
func initLogWithConfig(config LoggerConfig) *zap.SugaredLogger {
	var cores []zapcore.Core
	encoder := getEncoder()
	level := getLogLevel(config.LogLevel)

	// Add terminal output if needed
	if config.OutputMode == OutputTerminal || config.OutputMode == OutputBoth {
		terminalCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level)
		cores = append(cores, terminalCore)
	}

	// Add file output if needed
	if config.OutputMode == OutputFile || config.OutputMode == OutputBoth {
		fileCore := zapcore.NewCore(encoder, getLogWriter(config.LogDir), level)
		cores = append(cores, fileCore)
	}

	// If no valid output mode, default to terminal
	if len(cores) == 0 {
		terminalCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level)
		cores = append(cores, terminalCore)
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())
	sugarLogger := logger.Sugar()
	return sugarLogger
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

func getEncoder() zapcore.Encoder {
	loggerConfig := zap.NewProductionEncoderConfig()
	loggerConfig.TimeKey = "timestamp"
	loggerConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z07:00")
	loggerConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.FunctionKey = "func"
	return zapcore.NewJSONEncoder(loggerConfig)
}

func getLogWriter(logDir string) zapcore.WriteSyncer {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// If can't create directory, fallback to current directory
		logDir = "."
	}

	logFile := logDir + "/" + prefix() + ".log"
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	})
	return ws
}

// WithContext creates a new logger instance with context information.
// If the context contains a request ID, it will be automatically included in logs.
func (l Logger) WithContext(ctx context.Context) Logger {
	return Logger{
		log:          l.log,
		ctx:          ctx,
		level:        "",
		message:      "",
		data:         make([]any, 0),
		hasData:      false,
		requestIDKey: l.requestIDKey,
	}
}

// Debug sets the log level to debug and message.
func (l Logger) Debug(msg string) Logger {
	l.level = "debug"
	l.message = msg
	return l
}

// Info sets the log level to info and message.
func (l Logger) Info(msg string) Logger {
	l.level = "info"
	l.message = msg
	return l
}

// Warn sets the log level to warn and message.
func (l Logger) Warn(msg string) Logger {
	l.level = "warn"
	l.message = msg
	return l
}

// Error sets the log level to error and message.
func (l Logger) Error(msg string) Logger {
	l.level = "error"
	l.message = msg
	return l
}

// Fatal sets the log level to fatal and message.
func (l Logger) Fatal(msg string) Logger {
	l.level = "fatal"
	l.message = msg
	return l
}

// Panic sets the log level to panic and message.
func (l Logger) Panic(msg string) Logger {
	l.level = "panic"
	l.message = msg
	return l
}

// Data adds key-value pairs to the log data.
func (l Logger) Data(key string, value any) Logger {
	l.data = append(l.data, key, value)
	l.hasData = true
	return l
}

// ErrorData adds error information to the log data.
func (l Logger) ErrorData(err error) Logger {
	if err != nil {
		l.data = append(l.data, "error", err.Error())
		l.hasData = true
	}
	return l
}

// Send executes the log operation.
func (l Logger) Send() {
	requestID := GetRequestID(l.ctx)

	// Prepare log data
	logData := make([]any, 0, len(l.data)+2)
	if requestID != "" {
		logData = append(logData, l.requestIDKey, requestID)
	}
	logData = append(logData, l.data...)

	// Always use structured logging if we have any data (including request ID)
	hasStructuredData := len(logData) > 0

	// Log based on level
	switch l.level {
	case "debug":
		if hasStructuredData {
			l.log.Debugw(l.message, logData...)
		} else {
			l.log.Debug(l.message)
		}
	case "info":
		if hasStructuredData {
			l.log.Infow(l.message, logData...)
		} else {
			l.log.Info(l.message)
		}
	case "warn":
		if hasStructuredData {
			l.log.Warnw(l.message, logData...)
		} else {
			l.log.Warn(l.message)
		}
	case "error":
		if hasStructuredData {
			l.log.Errorw(l.message, logData...)
		} else {
			l.log.Error(l.message)
		}
	case "fatal":
		if hasStructuredData {
			l.log.Fatalw(l.message, logData...)
		} else {
			l.log.Fatal(l.message)
		}
	case "panic":
		if hasStructuredData {
			l.log.Panicw(l.message, logData...)
		} else {
			l.log.Panic(l.message)
		}
	}
}

// Close syncs all buffered logs and closes the logger.
// It ignores any sync errors as recommended by the underlying logger documentation.
func (l Logger) Close() {
	_ = l.log.Sync()
}
