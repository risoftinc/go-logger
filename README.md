# Gologger Package

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/go.risoftinc.com/gologger)](https://goreportcard.com/report/go.risoftinc.com/gologger)

A structured logging solution for Go applications using zap gologger. This package provides a simplified interface for logging with support for multiple output modes and log levels.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Basic Usage](#basic-usage)
  - [Custom Configuration](#custom-configuration)
  - [Logging with Data](#logging-with-data)
  - [Custom Request ID Key](#custom-request-id-key)
- [Context Support](#context-support)
- [Method Chaining Behavior](#method-chaining-behavior)
- [Configuration Options](#configuration-options)
- [API Reference](#api-reference)
- [Log File Configuration](#log-file-configuration)
- [Performance & Thread Safety](#performance--thread-safety)
- [Troubleshooting](#troubleshooting)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Multiple Output Modes**: Terminal, file, or both
- **Configurable Log Levels**: Debug, Info, Warn, Error
- **Structured Logging**: JSON format with timestamps and caller information
- **Caller Configuration**: Control whether to show caller information in logs (default: enabled)
- **Log Rotation**: Automatic log file rotation using lumberjack
- **Request Tracing**: Support for request ID tracking with custom key configuration (commonly used for HTTP request tracing)
- **Method Chaining**: Fluent API for clean, readable code
- **Context Support**: Automatic request ID inclusion from Go context
- **Thread Safe**: Built on zap's thread-safe foundation

## Installation

```bash
go get go.risoftinc.com/gologger
```

## Quick Start

### Basic Usage

```go
package main

import (
    "go.risoftinc.com/gologger"
)

func main() {
    // Create logger with default configuration
    log := gologger.NewLogger()
    defer log.Close()

    // Simple logging with method chaining
    log.Info("Application started").Send()
    log.Warn("This is a warning").Send()
    log.Error("An error occurred").Send()
}
```

### Custom Configuration

```go
package main

import (
    "go.risoftinc.com/gologger"
)

func main() {
    // Create logger with custom configuration
    config := gologger.LoggerConfig{
        OutputMode: gologger.OutputFile,
        LogLevel:   gologger.LevelInfo,
        LogDir:     "logs",
        ShowCaller: true, // Show caller information (default: true)
    }
    
    log := gologger.NewLoggerWithConfig(config)
    defer log.Close()

    log.Info("Application started with custom config").Send()
}
```

### Caller Configuration

```go
package main

import (
    "go.risoftinc.com/gologger"
)

func main() {
    // Logger with caller information (default behavior)
    logWithCaller := gologger.NewLogger()
    defer logWithCaller.Close()
    
    logWithCaller.Info("This log will show caller information").Send()

    // Logger without caller information
    configWithoutCaller := gologger.LoggerConfig{
        OutputMode: gologger.OutputTerminal,
        LogLevel:   gologger.LevelInfo,
        LogDir:     "logger",
        ShowCaller: false, // Disable caller information
    }
    
    logWithoutCaller := gologger.NewLoggerWithConfig(configWithoutCaller)
    defer logWithoutCaller.Close()
    
    logWithoutCaller.Info("This log will NOT show caller information").Send()
}
```

### Log Rotation Configuration

```go
package main

import (
    "go.risoftinc.com/gologger"
)

func main() {
    // Custom rotation configuration
    config := gologger.LoggerConfig{
        OutputMode: gologger.OutputBoth,
        LogLevel:   gologger.LevelInfo,
        LogDir:     "logs",
        LogRotation: &gologger.LogRotationConfig{
            MaxSize:    5,   // 5 MB instead of default 10 MB
            MaxBackups: 5,   // Keep 5 backup files instead of default 3
            MaxAge:     14,  // Keep files for 14 days instead of default 28
            Compress:   false, // Don't compress rotated files
        },
    }

    log := gologger.NewLoggerWithConfig(config)
    defer log.Close()

    log.Info("Custom rotation configuration example").Send()
}
```

### Logging with Data

```go
// Log with additional data using method chaining
log.Info("User login").
    Data("user_id", 123).
    Data("ip", "192.168.1.1").
    Send()

// Context-based logging (RECOMMENDED)
ctx := gologger.WithRequestID(context.Background(), "req-123")
log.WithContext(ctx).
    Info("API request").
    Data("endpoint", "/api/users").
    Data("method", "GET").
    Send()
```

### Custom Request ID Key

```go
// Create logger with custom request ID key
config := gologger.LoggerConfig{
    OutputMode:   gologger.OutputTerminal,
    LogLevel:     gologger.LevelInfo,
    RequestIDKey: "trace_id", // Custom key instead of "request-id"
}

log := gologger.NewLoggerWithConfig(config)
defer log.Close()

// Add request ID to context
ctx := gologger.WithRequestID(context.Background(), "trace-456")

// Log with custom key
log.WithContext(ctx).
    Info("Processing request").
    Data("user_id", 123).
    Send()
// Output: {"level":"info","msg":"Processing request","trace_id":"trace-456","user_id":123}
```

## Context Support

The logger supports Go's context package for request tracing with automatic request ID inclusion:

```go
// Add request ID to context
ctx := gologger.WithRequestID(context.Background(), "req-123")

// Get request ID from context
requestID := gologger.GetRequestID(ctx)

// Log with context (automatically includes request ID if present)
log.WithContext(ctx).
    Info("Processing request").
    Data("step", "validation").
    Send()

// Even without additional data, request ID will still be included
log.WithContext(ctx).Info("Simple message").Send()
```

### HTTP Request Flow Example

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Add request ID to context
    ctx := gologger.WithRequestID(r.Context(), generateRequestID())
    
    // All subsequent logs will include the request ID
    log.WithContext(ctx).
        Info("Request started").
        Data("method", r.Method).
        Data("path", r.URL.Path).
        Send()
    
    // Process request...
    processUser(ctx, log)
    
    log.WithContext(ctx).
        Info("Request completed").
        Data("status", 200).
        Send()
}

func processUser(ctx context.Context, log gologger.Logger) {
    // Error handling example
    err := validateUser()
    if err != nil {
        log.WithContext(ctx).
            Error("User validation failed").
            ErrorData(err).
            Data("step", "validation").
            Send()
        return
    }
    
    // Success logging
    log.WithContext(ctx).
        Info("User processed successfully").
        Data("user_id", 123).
        Data("duration", "150ms").
        Send()
}
```

### Echo Framework Example

Request ID is commonly used for tracing HTTP requests. Here's a simple example with Echo framework:

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/labstack/echo/v4"
    "go.risoftinc.com/gologger"
)

func main() {
    log := gologger.NewLogger()
    defer log.Close()
    
    e := echo.New()
    
    // Request ID middleware
    e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
            c.Set("request_id", requestID)
            
            ctx := gologger.WithRequestID(c.Request().Context(), requestID)
            c.SetRequest(c.Request().WithContext(ctx))
            
            c.Response().Header().Set("X-Request-ID", requestID)
            return next(c)
        }
    })
    
    // Logging middleware
    e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            start := time.Now()
            ctx := c.Request().Context()
            
            log.WithContext(ctx).
                Info("Request started").
                Data("method", c.Request().Method).
                Data("path", c.Request().URL.Path).
                Send()
            
            err := next(c)
            
            log.WithContext(ctx).
                Info("Request completed").
                Data("status", c.Response().Status).
                Data("duration", time.Since(start).String()).
                Send()
            
            return err
        }
    })
    
    e.GET("/api/users", func(c echo.Context) error {
        ctx := c.Request().Context()
        log.WithContext(ctx).Info("Getting users").Send()
        
        users := []map[string]interface{}{
            {"id": 1, "name": "John Doe"},
            {"id": 2, "name": "Jane Smith"},
        }
        
        return c.JSON(http.StatusOK, map[string]interface{}{
            "users": users,
            "request_id": c.Get("request_id"),
        })
    })
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

## Method Chaining Behavior

The new method chaining API provides a clean, fluent interface:

### With Request ID in Context

```go
ctx := gologger.WithRequestID(context.Background(), "req-123")

// No additional data - still includes request ID
log.WithContext(ctx).Info("Simple message").Send()
// Output: {"level":"info","msg":"Simple message","request-id":"req-123"}

// With additional data - includes both request ID and data
log.WithContext(ctx).
    Info("Processing").
    Data("step", "validation").
    Data("user_id", 123).
    Send()
// Output: {"level":"info","msg":"Processing","request-id":"req-123","step":"validation","user_id":123}
```

### Without Request ID in Context

```go
ctx := context.Background()

// No additional data - simple message
log.WithContext(ctx).Info("Simple message").Send()
// Output: {"level":"info","msg":"Simple message"}

// With additional data - includes data only
log.WithContext(ctx).
    Info("Processing").
    Data("step", "validation").
    Send()
// Output: {"level":"info","msg":"Processing","step":"validation"}
```

### With Custom Request ID Key

```go
// Create logger with custom request ID key
config := gologger.LoggerConfig{
    OutputMode:   gologger.OutputTerminal,
    LogLevel:     gologger.LevelInfo,
    RequestIDKey: "trace_id",
}
log := gologger.NewLoggerWithConfig(config)

ctx := gologger.WithRequestID(context.Background(), "trace-456")

// Uses custom key "trace_id" instead of "request-id"
log.WithContext(ctx).Info("Processing with trace ID").Send()
// Output: {"level":"info","msg":"Processing with trace ID","trace_id":"trace-456"}

// With additional data
log.WithContext(ctx).
    Info("Processing").
    Data("step", "validation").
    Data("user_id", 123).
    Send()
// Output: {"level":"info","msg":"Processing","trace_id":"trace-456","step":"validation","user_id":123}
```

### Error Handling

```go
err := errors.New("database connection failed")
log.WithContext(ctx).
    Error("Operation failed").
    ErrorData(err).
    Data("retry_count", 3).
    Send()
// Output: {"level":"error","msg":"Operation failed","request-id":"req-123","error":"database connection failed","retry_count":3}
```

### Key Benefits

1. **Fluent API**: Clean method chaining
2. **Automatic Request ID**: If present in context, request ID is automatically included
3. **Custom Request ID Key**: Configure custom keys for different services (trace_id, correlation_id, etc.)
4. **Flexible Data**: Add as many data fields as needed
5. **Error Handling**: Dedicated `ErrorData()` method for errors
6. **Type Safe**: Compile-time checking for method chaining

## API Reference

### Constructor Functions

- `NewLogger()`: Creates logger with default configuration
- `NewLoggerWithConfig(config gologger.LoggerConfig)`: Creates logger with custom configuration

### gologger.LoggerConfig Fields

- `OutputMode string`: Output mode (`OutputTerminal`, `OutputFile`, `OutputBoth`)
- `LogLevel string`: Log level (`LevelDebug`, `LevelInfo`, `LevelWarn`, `LevelError`)
- `LogDir string`: Directory for log files
- `RequestIDKey string`: Custom key for request ID in logs (default: `"request-id"`)
- `ShowCaller bool`: Whether to show caller information in logs (default: `true`)

### Context Functions

- `WithRequestID(ctx context.Context, requestID string) context.Context`: Adds request ID to context
- `GetRequestID(ctx context.Context) string`: Retrieves request ID from context

### Method Chaining API

#### Log Level Methods
- `Debug(msg string) gologger.Logger` - Sets debug level and message
- `Info(msg string) gologger.Logger` - Sets info level and message
- `Warn(msg string) gologger.Logger` - Sets warn level and message
- `Error(msg string) gologger.Logger` - Sets error level and message
- `Fatal(msg string) gologger.Logger` - Sets fatal level and message
- `Panic(msg string) gologger.Logger` - Sets panic level and message

#### Data Methods
- `Data(key string, value any) gologger.Logger` - Adds key-value pair to log data
- `ErrorData(err error) gologger.Logger` - Adds error information to log data

#### Context Methods
- `WithContext(ctx context.Context) gologger.Logger` - Creates logger with context

#### Execution Method
- `Send()` - Executes the log operation

### Utility Methods

- `Close()`: Syncs and closes the logger

## Configuration Options

### gologger.LoggerConfig Structure

```go
type gologger.LoggerConfig struct {
    OutputMode    string              // Output mode: OutputTerminal, OutputFile, or OutputBoth
    LogLevel      string              // Log level: LevelDebug, LevelInfo, LevelWarn, or LevelError
    LogDir        string              // Directory for log files
    RequestIDKey  string              // Custom key for request ID in logs (default: "request-id")
    ShowCaller    bool                // Whether to show caller information in logs (default: true)
    LogRotation   *LogRotationConfig  // Log rotation configuration (optional, uses defaults if nil)
}

type gologger.LogRotationConfig struct {
    MaxSize    int  // Maximum size in megabytes before rotation (default: 10)
    MaxBackups int  // Maximum number of old log files to retain (default: 3)
    MaxAge     int  // Maximum number of days to retain old log files (default: 28)
    Compress   bool // Whether to compress rotated log files (default: true)
}
```

### Custom Request ID Key

You can customize the key used for request ID in logs:

```go
// Default request ID key
log := gologger.NewLogger()
// Output: {"level":"info","msg":"Message","request-id":"req-123"}

// Custom request ID key
config := gologger.LoggerConfig{
    OutputMode:   gologger.OutputTerminal,
    LogLevel:     gologger.LevelInfo,
    RequestIDKey: "trace_id",
}
log := gologger.NewLoggerWithConfig(config)
// Output: {"level":"info","msg":"Message","trace_id":"req-123"}

// Different keys for different services
apiConfig := gologger.LoggerConfig{
    RequestIDKey: "request_id",
}
apiLogger := gologger.NewLoggerWithConfig(apiConfig)

bgConfig := gologger.LoggerConfig{
    RequestIDKey: "job_id",
}
bgLogger := gologger.NewLoggerWithConfig(bgConfig)
```

### Common Request ID Key Patterns

- `"request-id"` - Default, good for HTTP APIs
- `"trace_id"` - For distributed tracing
- `"correlation_id"` - For business correlation
- `"session_id"` - For user sessions
- `"job_id"` - For background jobs
- `"transaction_id"` - For database transactions

## Log File Configuration

### Default Rotation Settings

Log files are automatically rotated with the following default settings:
- Maximum file size: 10 MB
- Maximum backup files: 3
- Maximum age: 28 days
- Compression: Enabled

Log files are named with the pattern: `logger-YYYY-MM-DD.log`

### Custom Rotation Configuration

You can customize log rotation settings by providing a `LogRotationConfig`:

```go
// Custom rotation configuration
config := gologger.LoggerConfig{
    OutputMode: gologger.OutputBoth,
    LogLevel:   gologger.LevelInfo,
    LogDir:     "logs",
    LogRotation: &gologger.LogRotationConfig{
        MaxSize:    5,   // 5 MB instead of default 10 MB
        MaxBackups: 5,   // Keep 5 backup files instead of default 3
        MaxAge:     14,  // Keep files for 14 days instead of default 28
        Compress:   false, // Don't compress rotated files
    },
}

log := gologger.NewLoggerWithConfig(config)
```

### Partial Configuration

You can set only specific rotation parameters, and the rest will use defaults:

```go
// Only set MaxSize, others will use defaults
config := gologger.LoggerConfig{
    OutputMode: gologger.OutputFile,
    LogLevel:   gologger.LevelDebug,
    LogDir:     "logs",
    LogRotation: &gologger.LogRotationConfig{
        MaxSize: 20, // Only this is set, others use defaults
    },
}
```

### Using Default Rotation

If you don't specify `LogRotation` or set it to `nil`, all rotation settings will use the defaults:

```go
// Uses all default rotation settings
config := gologger.LoggerConfig{
    OutputMode: gologger.OutputBoth,
    LogLevel:   gologger.LevelInfo,
    LogDir:     "logs",
    // LogRotation is nil, so defaults are used
}
```

## Performance & Thread Safety

### Performance
- **Zero-allocation logging**: Built on zap's high-performance foundation
- **Structured logging**: JSON output for easy parsing and analysis
- **Method chaining**: Minimal overhead, returns new instances for immutability

### Thread Safety
- **Concurrent safe**: All logger methods are thread-safe
- **Immutable design**: Method chaining returns new instances
- **Context propagation**: Safe to pass logger instances across goroutines

### Memory Usage
- **Efficient**: Uses zap's object pool for reduced GC pressure
- **Configurable**: Adjust log levels to control verbosity
- **Rotation**: Automatic log file rotation prevents disk space issues

## Troubleshooting

### Common Issues

#### 1. Log Files Not Created
```go
// Ensure log directory exists and has write permissions
config := gologger.LoggerConfig{
    OutputMode: gologger.OutputFile,
    LogDir:     "logs", // Make sure this directory exists
}
```

#### 2. Request ID Not Appearing
```go
// Make sure to use WithContext() method
ctx := gologger.WithRequestID(context.Background(), "req-123")
log.WithContext(ctx).Info("Message").Send() // ✅ Correct

// Don't forget to call Send()
log.WithContext(ctx).Info("Message") // ❌ Won't log
```

#### 3. Custom Request ID Key Not Working
```go
// Ensure RequestIDKey is set in config
config := gologger.LoggerConfig{
    RequestIDKey: "trace_id", // Must be set
}
log := gologger.NewLoggerWithConfig(config)
```

#### 4. Method Chaining Not Working
```go
// Each method returns a new gologger.Logger instance
logger := log.Info("Message").Data("key", "value") // Returns new instance
logger.Send() // Must call Send() on the returned instance
```

### Debug Mode
```go
// Enable debug logging to see internal operations
config := gologger.LoggerConfig{
    LogLevel: gologger.LevelDebug,
}
log := gologger.NewLoggerWithConfig(config)
```

## Dependencies

- [go.uber.org/zap](https://github.com/uber-go/zap): High-performance structured logging
- [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack): Log rotation

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 Risoftinc.
