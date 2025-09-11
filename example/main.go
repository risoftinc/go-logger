package main

import (
	"context"
	"errors"

	"go.risoftinc.com/gologger"
)

func main() {
	// Example 1: Basic logging with method chaining
	log := gologger.NewLogger()
	defer log.Close()

	// Simple logging
	log.Info("Application started").Send()
	log.Debug("Debug message").Send()
	log.Warn("Warning message").Send()
	log.Error("Error message").Send()

	// Example 2: Logging with data
	log.Info("User action").
		Data("user_id", 123).
		Data("action", "login").
		Data("ip", "192.168.1.1").
		Send()

	// Example 3: Context-based logging (RECOMMENDED WAY)
	ctx := context.Background()
	ctx = gologger.WithRequestID(ctx, "req-67890")

	// Log with context - request ID will be automatically included
	log.WithContext(ctx).
		Info("API request").
		Data("endpoint", "/api/users").
		Data("method", "GET").
		Data("status", 200).
		Send()

	// Example 4: Error logging with error data
	err := errors.New("database connection failed")
	log.WithContext(ctx).
		Error("Database operation failed").
		ErrorData(err).
		Data("retry_count", 3).
		Data("timeout", "5s").
		Send()

	// Example 5: Multiple data fields
	log.WithContext(ctx).
		Debug("Processing request").
		Data("step", "validation").
		Data("component", "auth").
		Data("duration", "150ms").
		Send()

	// Example 6: Context without request ID
	ctxNoID := context.Background()
	log.WithContext(ctxNoID).
		Info("No request ID in context").
		Data("component", "system").
		Send()

	// Example 7: Simulating HTTP request flow
	simulateHTTPRequest(log)
}

func simulateHTTPRequest(log gologger.Logger) {
	// Simulate incoming HTTP request
	ctx := context.Background()
	ctx = gologger.WithRequestID(ctx, "req-http-001")

	log.WithContext(ctx).
		Info("HTTP request received").
		Data("method", "POST").
		Data("path", "/api/users").
		Send()

	// Simulate processing
	processUser(ctx, log)

	log.WithContext(ctx).
		Info("HTTP request completed").
		Data("status", 201).
		Send()
}

func processUser(ctx context.Context, log gologger.Logger) {
	log.WithContext(ctx).
		Info("Processing user data").
		Data("step", "validation").
		Send()

	// Simulate some processing
	log.WithContext(ctx).
		Debug("User validation passed").
		Data("rules", "email,password").
		Send()

	// Simulate database operation
	log.WithContext(ctx).
		Info("Saving user to database").
		Data("table", "users").
		Send()

	// Simulate potential error
	log.WithContext(ctx).
		Warn("Database query slow").
		Data("duration", "1.2s").
		Data("threshold", "1.0s").
		Send()

	// Simulate error handling
	err := errors.New("connection timeout")
	log.WithContext(ctx).
		Error("Database operation failed").
		ErrorData(err).
		Data("retry_count", 3).
		Send()
}
