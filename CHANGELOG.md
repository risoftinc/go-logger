# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- 

### Fixed
- 

### Features
- 

## [1.3.0] - 2025-09-12

### Added
- **Caller Configuration**: Added `ShowCaller` field to `LoggerConfig` struct to control whether caller information is displayed in logs
- **Default Caller Behavior**: Caller information is now shown by default (`ShowCaller: true`)
- **Caller Control Examples**: Updated examples to demonstrate caller configuration usage

### Changed
- **LoggerConfig Structure**: Added `ShowCaller bool` field to `LoggerConfig` struct
- **Logger Struct**: Added `showCaller bool` field to `Logger` struct to store caller configuration
- **initLogWithConfig Function**: Modified to conditionally add caller information based on `ShowCaller` configuration
- **WithContext Method**: Updated to preserve caller configuration when creating new logger instances

### Technical Details
- Caller information is controlled by `zap.AddCaller()` and `zap.AddCallerSkip(1)` options
- When `ShowCaller: false`, these options are not applied to the zap logger
- Default behavior maintains backward compatibility by showing caller information
- All existing functionality remains unchanged

### Usage Example
```go
// Logger with caller information (default)
log := gologger.NewLogger()

// Logger without caller information
config := gologger.LoggerConfig{
    OutputMode: gologger.OutputTerminal,
    LogLevel:   gologger.LevelInfo,
    LogDir:     "logger",
    ShowCaller: false, // Disable caller information
}
log := gologger.NewLoggerWithConfig(config)
