package main

import (
	"context"

	"github.com/risoftinc/gologger"
)

func ExampleLogRotation() {
	// Example 1: Using default rotation settings
	log := gologger.NewLogger()
	defer log.Close()

	log.Info("Using default rotation settings").Send()

	// Example 2: Custom rotation configuration
	customConfig := gologger.LoggerConfig{
		OutputMode: gologger.OutputBoth,
		LogLevel:   gologger.LevelInfo,
		LogDir:     "custom_logs",
		LogRotation: &gologger.LogRotationConfig{
			MaxSize:    5,     // 5 MB instead of default 10 MB
			MaxBackups: 5,     // Keep 5 backup files instead of default 3
			MaxAge:     14,    // Keep files for 14 days instead of default 28
			Compress:   false, // Don't compress rotated files
		},
	}

	customLog := gologger.NewLoggerWithConfig(customConfig)
	defer customLog.Close()

	customLog.Info("Using custom rotation settings").Send()

	// Example 3: Partial rotation configuration (only some fields set)
	partialConfig := gologger.LoggerConfig{
		OutputMode: gologger.OutputFile,
		LogLevel:   gologger.LevelDebug,
		LogDir:     "partial_logs",
		LogRotation: &gologger.LogRotationConfig{
			MaxSize: 20, // Only set MaxSize, others will use defaults
			// MaxBackups will default to 3
			// MaxAge will default to 28
			// Compress will default to true
		},
	}

	partialLog := gologger.NewLoggerWithConfig(partialConfig)
	defer partialLog.Close()

	partialLog.Info("Using partial rotation settings").Send()

	// Example 4: No rotation configuration (uses all defaults)
	defaultRotationConfig := gologger.LoggerConfig{
		OutputMode: gologger.OutputBoth,
		LogLevel:   gologger.LevelWarn,
		LogDir:     "default_rotation_logs",
		// LogRotation is nil, so all rotation settings will use defaults
	}

	defaultRotationLog := gologger.NewLoggerWithConfig(defaultRotationConfig)
	defer defaultRotationLog.Close()

	defaultRotationLog.Info("Using default rotation settings (nil config)").Send()

	// Example 5: Context-based logging with custom rotation
	ctx := context.Background()
	ctx = gologger.WithRequestID(ctx, "rotation-demo-001")

	customLog.WithContext(ctx).
		Info("Log rotation demo with context").
		Data("rotation_config", "custom").
		Data("max_size_mb", 5).
		Data("max_backups", 5).
		Data("max_age_days", 14).
		Data("compress", false).
		Send()

	// Example 6: Demonstrating different log levels with rotation
	levels := []string{"debug", "info", "warn", "error"}
	for i, level := range levels {
		customLog.WithContext(ctx).
			Info("Testing log rotation").
			Data("level", level).
			Data("iteration", i+1).
			Data("message", "This is a test message to demonstrate log rotation").
			Send()
	}
}
