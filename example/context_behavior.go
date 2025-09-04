package main

import (
	"context"

	"github.com/risoftinc/gologger"
)

// ContextBehaviorExample demonstrates different scenarios when using context-aware logging
func ContextBehaviorExample() {
	log := gologger.NewLogger()
	defer log.Close()

	// Scenario 1: Context with request ID, no additional data
	ctx := gologger.WithRequestID(context.Background(), "req-001")

	log.WithContext(ctx).Info("Simple message with request ID only").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Simple message with request ID only","request-id":"req-001"}

	// Scenario 2: Context with request ID, with additional data
	log.WithContext(ctx).
		Info("Message with request ID and data").
		Data("user_id", 123).
		Data("action", "login").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"Message with request ID and data","request-id":"req-001","user_id":123,"action":"login"}

	// Scenario 3: Context without request ID, no additional data
	ctxNoID := context.Background()
	log.WithContext(ctxNoID).Info("Simple message without request ID").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Simple message without request ID"}

	// Scenario 4: Context without request ID, with additional data
	log.WithContext(ctxNoID).
		Info("Message without request ID but with data").
		Data("component", "auth").
		Data("status", "success").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"Message without request ID but with data","component":"auth","status":"success"}

	// Scenario 5: Different log levels with context
	log.WithContext(ctx).Debug("Debug message with request ID").Send()
	log.WithContext(ctx).Warn("Warning message with request ID").Send()
	log.WithContext(ctx).Error("Error message with request ID").Send()

	// Scenario 6: Empty data slice (same as no data)
	log.WithContext(ctx).Info("Message with empty data slice").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Message with empty data slice","request-id":"req-001"}

	// Scenario 7: Nil data (demonstrates that nil data is handled gracefully)
	log.WithContext(ctx).Info("Message with nil data").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Message with nil data","request-id":"req-001"}
}

// DemonstrateHTTPFlow shows how it works in a typical HTTP request flow
func DemonstrateHTTPFlow() {
	log := gologger.NewLogger()
	defer log.Close()

	// Simulate HTTP request processing
	ctx := gologger.WithRequestID(context.Background(), "req-http-123")

	// Step 1: Request received (no additional data)
	log.WithContext(ctx).Info("HTTP request received").Send()
	// Output: {"level":"info","timestamp":"...","msg":"HTTP request received","request-id":"req-http-123"}

	// Step 2: Processing with some data
	log.WithContext(ctx).
		Info("Processing request").
		Data("method", "POST").
		Data("path", "/api/users").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"Processing request","request-id":"req-http-123","method":"POST","path":"/api/users"}

	// Step 3: Database operation (no additional data)
	log.WithContext(ctx).Info("Saving to database").Send()
	// Output: {"level":"info","timestamp":"...","msg":"Saving to database","request-id":"req-http-123"}

	// Step 4: Error handling (with error data)
	log.WithContext(ctx).
		Error("Database error occurred").
		Data("error", "connection timeout").
		Data("retry_count", 3).
		Send()
	// Output: {"level":"error","timestamp":"...","msg":"Database error occurred","request-id":"req-http-123","error":"connection timeout","retry_count":3}

	// Step 5: Request completed (minimal data)
	log.WithContext(ctx).
		Info("Request completed").
		Data("status", 500).
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"Request completed","request-id":"req-http-123","status":500}
}

// ShowComparison demonstrates the difference between old and new approaches
func ShowComparison() {
	log := gologger.NewLogger()
	defer log.Close()

	requestID := "req-comparison-456"
	ctx := gologger.WithRequestID(context.Background(), requestID)

	// OLD WAY: Manual request ID passing (method chaining approach)
	log.WithContext(ctx).
		Info("Old way - manual request ID").
		Data("step", "validation").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"Old way - manual request ID","request-id":"req-comparison-456","step":"validation"}

	// NEW WAY: Context-aware (cleaner)
	log.WithContext(ctx).
		Info("New way - context aware").
		Data("step", "validation").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"New way - context aware","request-id":"req-comparison-456","step":"validation"}

	// OLD WAY: Without request ID
	log.Info("Old way - no request ID").
		Data("step", "validation").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"Old way - no request ID","data":["step","validation"]}

	// NEW WAY: Without request ID (context doesn't have it)
	ctxNoID := context.Background()
	log.WithContext(ctxNoID).
		Info("New way - no request ID").
		Data("step", "validation").
		Send()
	// Output: {"level":"info","timestamp":"...","msg":"New way - no request ID","step":"validation"}
}
