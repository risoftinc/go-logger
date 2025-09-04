# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

## [1.0.0] - 2025-09-04

### Added
- Initial release of logger package
- Method chaining API similar
- Context-aware logging with automatic request ID inclusion
- Custom request ID key configuration for different services
- Multiple output modes (terminal, file, both)
- Configurable log levels (debug, info, warn, error)
- Structured JSON logging with zap
- Log rotation with lumberjack
- Error handling with dedicated ErrorData() method
- Comprehensive documentation and examples

### Features
- **Method Chaining**: Clean fluent API
  ```go
  log.WithContext(ctx).
      Info("Processing request").
      Data("user_id", 123).
      Data("step", "validation").
      Send()
  ```

- **Context Support**: Automatic request ID inclusion
  ```go
  ctx := logger.WithRequestID(context.Background(), "req-123")
  log.WithContext(ctx).Info("Request started").Send()
  ```

- **Custom Request ID Key**: Configure different keys for different services
  ```go
  // API Service
  apiConfig := logger.LoggerConfig{
      RequestIDKey: "request_id",
  }
  apiLogger := logger.NewLoggerWithConfig(apiConfig)
  
  // Background Service
  bgConfig := logger.LoggerConfig{
      RequestIDKey: "job_id",
  }
  bgLogger := logger.NewLoggerWithConfig(bgConfig)
  ```

- **Flexible Configuration**: Multiple output modes and log levels
  ```go
  config := logger.LoggerConfig{
      OutputMode: logger.OutputBoth,
      LogLevel:   logger.LevelInfo,
      LogDir:     "logs",
  }
  log := logger.NewLoggerWithConfig(config)
  ```

- **Error Handling**: Dedicated error logging
  ```go
  log.WithContext(ctx).
      Error("Operation failed").
      ErrorData(err).
      Data("retry_count", 3).
      Send()
  ```

- **Log Rotation**: Automatic file rotation
  - Max file size: 10 MB
  - Max backup files: 3
  - Max age: 28 days
  - Compression: Enabled

### Technical Details
- Built on top of [go.uber.org/zap](https://github.com/uber-go/zap)
- Uses [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack) for log rotation
- Thread-safe logging
- Zero-allocation logging for high performance
- JSON structured output for easy parsing

### API Reference
- `NewLogger()` - Create logger with default config
- `NewLoggerWithConfig(config LoggerConfig)` - Create logger with custom config
- `WithContext(ctx context.Context) Logger` - Create logger with context
- `WithRequestID(ctx context.Context, requestID string) context.Context` - Add request ID to context
- `GetRequestID(ctx context.Context) string` - Get request ID from context
- `Debug(msg string) Logger` - Set debug level
- `Info(msg string) Logger` - Set info level
- `Warn(msg string) Logger` - Set warn level
- `Error(msg string) Logger` - Set error level
- `Fatal(msg string) Logger` - Set fatal level
- `Panic(msg string) Logger` - Set panic level
- `Data(key string, value any) Logger` - Add key-value data
- `ErrorData(err error) Logger` - Add error information
- `Send()` - Execute the log operation
- `Close()` - Close and sync logger

### Examples
- Basic logging examples in `example/main.go`
- HTTP handler examples in `example/http_example.go`
- Context behavior examples in `example/context_behavior.go`
