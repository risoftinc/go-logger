package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/risoftinc/gologger"
)

func customRequestIDExample() {
	fmt.Println("🔧 Custom Request ID Key Examples")
	fmt.Println(strings.Repeat("=", 50))

	// Example 1: Default request ID key
	fmt.Println("\n1. Default Request ID Key (request-id):")
	log1 := gologger.NewLogger()
	defer log1.Close()

	ctx1 := gologger.WithRequestID(context.Background(), "req-123")
	log1.WithContext(ctx1).Info("Using default request ID key").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Using default request ID key","request-id":"req-123"}

	// Example 2: Custom request ID key - "trace_id"
	fmt.Println("\n2. Custom Request ID Key (trace_id):")
	config2 := gologger.LoggerConfig{
		OutputMode:   gologger.OutputTerminal,
		LogLevel:     gologger.LevelInfo,
		LogDir:       "logs",
		RequestIDKey: "trace_id",
	}
	log2 := gologger.NewLoggerWithConfig(config2)
	defer log2.Close()

	ctx2 := gologger.WithRequestID(context.Background(), "trace-456")
	log2.WithContext(ctx2).Info("Using custom trace_id key").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Using custom trace_id key","trace_id":"trace-456"}

	// Example 3: Custom request ID key - "correlation_id"
	fmt.Println("\n3. Custom Request ID Key (correlation_id):")
	config3 := gologger.LoggerConfig{
		OutputMode:   gologger.OutputTerminal,
		LogLevel:     gologger.LevelInfo,
		LogDir:       "logs",
		RequestIDKey: "correlation_id",
	}
	log3 := gologger.NewLoggerWithConfig(config3)
	defer log3.Close()

	ctx3 := gologger.WithRequestID(context.Background(), "corr-789")
	log3.WithContext(ctx3).
		Info("Using custom correlation_id key").
		Data("user_id", 123).
		Data("action", "login").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"Using custom correlation_id key","correlation_id":"corr-789","user_id":123,"action":"login"}

	// Example 4: Custom request ID key - "session_id"
	fmt.Println("\n4. Custom Request ID Key (session_id):")
	config4 := gologger.LoggerConfig{
		OutputMode:   gologger.OutputTerminal,
		LogLevel:     gologger.LevelInfo,
		LogDir:       "logs",
		RequestIDKey: "session_id",
	}
	log4 := gologger.NewLoggerWithConfig(config4)
	defer log4.Close()

	ctx4 := gologger.WithRequestID(context.Background(), "sess-abc123")
	log4.WithContext(ctx4).
		Error("Database connection failed").
		Data("error", "connection timeout").
		Data("retry_count", 3).
		Send()
	// Output: {"level":"error","timestamp":"...","msg":"Database connection failed","session_id":"sess-abc123","error":"connection timeout","retry_count":3}

	// Example 5: Different keys for different services
	fmt.Println("\n5. Different Keys for Different Services:")

	// API Service using "request_id"
	apiConfig := gologger.LoggerConfig{
		OutputMode:   gologger.OutputTerminal,
		LogLevel:     gologger.LevelInfo,
		RequestIDKey: "request_id",
	}
	apiLogger := gologger.NewLoggerWithConfig(apiConfig)
	defer apiLogger.Close()

	// Background Service using "job_id"
	bgConfig := gologger.LoggerConfig{
		OutputMode:   gologger.OutputTerminal,
		LogLevel:     gologger.LevelInfo,
		RequestIDKey: "job_id",
	}
	bgLogger := gologger.NewLoggerWithConfig(bgConfig)
	defer bgLogger.Close()

	// API Service logging
	apiCtx := gologger.WithRequestID(context.Background(), "api-req-001")
	apiLogger.WithContext(apiCtx).Info("API request processed").Send()
	// Output: {"level":"info","timestamp":"...","msg":"API request processed","request_id":"api-req-001"}

	// Background Service logging
	bgCtx := gologger.WithRequestID(context.Background(), "job-bg-002")
	bgLogger.WithContext(bgCtx).Info("Background job completed").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Background job completed","job_id":"job-bg-002"}

	// Example 6: Empty request ID key (no request ID logging)
	fmt.Println("\n6. Empty Request ID Key (no request ID logging):")
	config6 := gologger.LoggerConfig{
		OutputMode:   gologger.OutputTerminal,
		LogLevel:     gologger.LevelInfo,
		RequestIDKey: "", // Empty key means no request ID logging
	}
	log6 := gologger.NewLoggerWithConfig(config6)
	defer log6.Close()

	ctx6 := gologger.WithRequestID(context.Background(), "req-999")
	log6.WithContext(ctx6).
		Info("This won't show request ID").
		Data("user_id", 456).
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"This won't show request ID","user_id":456}

	fmt.Println("\n✅ Custom Request ID Key Examples Completed!")
}
