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

## [1.1.0] - 2025-09-09

### Added
- **Log Rotation Configuration**: Added configurable log file rotation settings
  - New `LogRotationConfig` struct with `MaxSize`, `MaxBackups`, `MaxAge`, and `Compress` fields
  - Added `LogRotation` field to `LoggerConfig` for optional rotation configuration
  - Default values maintained when rotation config is not provided (10MB max size, 3 backups, 28 days retention, compression enabled)
  - Support for partial configuration - only set specific fields, others use defaults
  - Updated `getLogWriter` function to use rotation configuration with fallback to defaults
  - Added comprehensive documentation and examples for log rotation configuration
  - Created `rotation_example.go` demonstrating various rotation configuration scenarios

### Changed
- Modified `getLogWriter` function signature to accept `LogRotationConfig` parameter
- Updated `initLogWithConfig` to pass rotation configuration to `getLogWriter`
- Enhanced README with detailed log rotation configuration documentation
