package gologger

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	log := NewLogger()
	defer log.Close()

	if log.log == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	config := LoggerConfig{
		OutputMode: OutputTerminal,
		LogLevel:   LevelInfo,
		LogDir:     "test_logs",
	}

	log := NewLoggerWithConfig(config)
	defer log.Close()

	if log.log == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestWithRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-123"

	ctxWithID := WithRequestID(ctx, requestID)
	retrievedID := GetRequestID(ctxWithID)

	if retrievedID != requestID {
		t.Errorf("Expected request ID %s, got %s", requestID, retrievedID)
	}
}

func TestGetRequestID_NoID(t *testing.T) {
	ctx := context.Background()
	requestID := GetRequestID(ctx)

	if requestID != "" {
		t.Errorf("Expected empty request ID, got %s", requestID)
	}
}

func TestWithContext(t *testing.T) {
	log := NewLogger()
	defer log.Close()

	ctx := context.Background()
	ctxWithID := WithRequestID(ctx, "test-123")

	contextLogger := log.WithContext(ctxWithID)

	if contextLogger.ctx == nil {
		t.Error("Expected context to be set")
	}

	// Test that context has request ID
	requestID := GetRequestID(contextLogger.ctx)
	if requestID != "test-123" {
		t.Errorf("Expected request ID test-123, got %s", requestID)
	}
}

func TestLogLevelMethods(t *testing.T) {
	log := NewLogger()
	defer log.Close()

	// Test that all level methods return logger instance
	debugLogger := log.Debug("debug message")
	if debugLogger.level != "debug" {
		t.Errorf("Expected level debug, got %s", debugLogger.level)
	}
	if debugLogger.message != "debug message" {
		t.Errorf("Expected message 'debug message', got %s", debugLogger.message)
	}

	infoLogger := log.Info("info message")
	if infoLogger.level != "info" {
		t.Errorf("Expected level info, got %s", infoLogger.level)
	}

	warnLogger := log.Warn("warn message")
	if warnLogger.level != "warn" {
		t.Errorf("Expected level warn, got %s", warnLogger.level)
	}

	errorLogger := log.Error("error message")
	if errorLogger.level != "error" {
		t.Errorf("Expected level error, got %s", errorLogger.level)
	}

	fatalLogger := log.Fatal("fatal message")
	if fatalLogger.level != "fatal" {
		t.Errorf("Expected level fatal, got %s", fatalLogger.level)
	}

	panicLogger := log.Panic("panic message")
	if panicLogger.level != "panic" {
		t.Errorf("Expected level panic, got %s", panicLogger.level)
	}
}

func TestDataMethod(t *testing.T) {
	log := NewLogger()
	defer log.Close()

	// Test single data addition
	loggerWithData := log.Info("test message").Data("key1", "value1")
	if len(loggerWithData.data) != 2 {
		t.Errorf("Expected 2 data items, got %d", len(loggerWithData.data))
	}
	if loggerWithData.data[0] != "key1" || loggerWithData.data[1] != "value1" {
		t.Errorf("Expected data ['key1', 'value1'], got %v", loggerWithData.data)
	}
	if !loggerWithData.hasData {
		t.Error("Expected hasData to be true")
	}

	// Test multiple data additions
	loggerWithMoreData := loggerWithData.Data("key2", "value2").Data("key3", 123)
	if len(loggerWithMoreData.data) != 6 {
		t.Errorf("Expected 6 data items, got %d", len(loggerWithMoreData.data))
	}
}

func TestErrorDataMethod(t *testing.T) {
	log := NewLogger()
	defer log.Close()

	err := errors.New("test error")
	loggerWithError := log.Info("test message").ErrorData(err)

	if len(loggerWithError.data) != 2 {
		t.Errorf("Expected 2 data items, got %d", len(loggerWithError.data))
	}
	if loggerWithError.data[0] != "error" {
		t.Errorf("Expected first data item to be 'error', got %s", loggerWithError.data[0])
	}
	if loggerWithError.data[1] != "test error" {
		t.Errorf("Expected second data item to be 'test error', got %s", loggerWithError.data[1])
	}
	if !loggerWithError.hasData {
		t.Error("Expected hasData to be true")
	}
}

func TestErrorDataMethod_NilError(t *testing.T) {
	log := NewLogger()
	defer log.Close()

	// Test with nil error
	loggerWithNilError := log.Info("test message").ErrorData(nil)

	if len(loggerWithNilError.data) != 0 {
		t.Errorf("Expected 0 data items for nil error, got %d", len(loggerWithNilError.data))
	}
	if loggerWithNilError.hasData {
		t.Error("Expected hasData to be false for nil error")
	}
}

func TestMethodChaining(t *testing.T) {
	log := NewLogger()
	defer log.Close()

	ctx := WithRequestID(context.Background(), "test-request")

	// Test complex method chaining
	logger := log.WithContext(ctx).
		Info("Processing request").
		Data("user_id", 123).
		Data("action", "login").
		Data("ip", "192.168.1.1")

	// Verify all data is set correctly
	if logger.level != "info" {
		t.Errorf("Expected level info, got %s", logger.level)
	}
	if logger.message != "Processing request" {
		t.Errorf("Expected message 'Processing request', got %s", logger.message)
	}
	if len(logger.data) != 6 {
		t.Errorf("Expected 6 data items, got %d", len(logger.data))
	}
	if !logger.hasData {
		t.Error("Expected hasData to be true")
	}

	// Verify context is preserved
	requestID := GetRequestID(logger.ctx)
	if requestID != "test-request" {
		t.Errorf("Expected request ID 'test-request', got %s", requestID)
	}
}

func TestSendMethod(t *testing.T) {
	// Create a temporary log file for testing
	tempDir := "test_logs"
	defer os.RemoveAll(tempDir)

	config := LoggerConfig{
		OutputMode: OutputFile,
		LogLevel:   LevelInfo,
		LogDir:     tempDir,
	}

	log := NewLoggerWithConfig(config)
	defer log.Close()

	// Test simple message
	log.Info("Test message").Send()

	// Test with data
	log.Info("Test with data").
		Data("key1", "value1").
		Data("key2", 123).
		Send()

	// Test with context
	ctx := WithRequestID(context.Background(), "test-request")
	log.WithContext(ctx).
		Info("Test with context").
		Data("user_id", 456).
		Send()

	// Give some time for file to be written
	time.Sleep(100 * time.Millisecond)

	// Check if log file was created
	logFile := tempDir + "/" + prefix() + ".log"
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Expected log file to be created at %s", logFile)
	}
}

func TestConstants(t *testing.T) {
	// Test output mode constants
	if OutputTerminal != "terminal" {
		t.Errorf("Expected OutputTerminal to be 'terminal', got %s", OutputTerminal)
	}
	if OutputFile != "file" {
		t.Errorf("Expected OutputFile to be 'file', got %s", OutputFile)
	}
	if OutputBoth != "both" {
		t.Errorf("Expected OutputBoth to be 'both', got %s", OutputBoth)
	}

	// Test log level constants
	if LevelDebug != "debug" {
		t.Errorf("Expected LevelDebug to be 'debug', got %s", LevelDebug)
	}
	if LevelInfo != "info" {
		t.Errorf("Expected LevelInfo to be 'info', got %s", LevelInfo)
	}
	if LevelWarn != "warn" {
		t.Errorf("Expected LevelWarn to be 'warn', got %s", LevelWarn)
	}
	if LevelError != "error" {
		t.Errorf("Expected LevelError to be 'error', got %s", LevelError)
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected zapcore.Level
	}{
		{LevelDebug, zapcore.DebugLevel},
		{LevelInfo, zapcore.InfoLevel},
		{LevelWarn, zapcore.WarnLevel},
		{LevelError, zapcore.ErrorLevel},
		{"invalid", zapcore.DebugLevel}, // default case
	}

	for _, test := range tests {
		result := getLogLevel(test.input)
		if result != test.expected {
			t.Errorf("getLogLevel(%s): expected %v, got %v", test.input, test.expected, result)
		}
	}
}

func TestCustomRequestIDKey(t *testing.T) {
	// Test with custom request ID key
	config := LoggerConfig{
		OutputMode:   OutputTerminal,
		LogLevel:     LevelInfo,
		LogDir:       "test_logs",
		RequestIDKey: "trace_id",
	}

	log := NewLoggerWithConfig(config)
	defer log.Close()

	if log.requestIDKey != "trace_id" {
		t.Errorf("Expected requestIDKey to be 'trace_id', got %s", log.requestIDKey)
	}

	// Test WithContext preserves custom request ID key
	ctx := WithRequestID(context.Background(), "trace-123")
	contextLogger := log.WithContext(ctx)

	if contextLogger.requestIDKey != "trace_id" {
		t.Errorf("Expected contextLogger.requestIDKey to be 'trace_id', got %s", contextLogger.requestIDKey)
	}
}

func TestDefaultRequestIDKey(t *testing.T) {
	// Test with default request ID key
	log := NewLogger()
	defer log.Close()

	if log.requestIDKey != "request-id" {
		t.Errorf("Expected requestIDKey to be 'request-id', got %s", log.requestIDKey)
	}
}

func TestEmptyRequestIDKey(t *testing.T) {
	// Test with empty request ID key
	config := LoggerConfig{
		OutputMode:   OutputTerminal,
		LogLevel:     LevelInfo,
		LogDir:       "test_logs",
		RequestIDKey: "",
	}

	log := NewLoggerWithConfig(config)
	defer log.Close()

	if log.requestIDKey != "request-id" {
		t.Errorf("Expected requestIDKey to be 'request-id' (default), got %s", log.requestIDKey)
	}
}

func TestPrefix(t *testing.T) {
	prefix := prefix()

	// Check if prefix starts with "logger-"
	if !strings.HasPrefix(prefix, "logger-") {
		t.Errorf("Expected prefix to start with 'logger-', got %s", prefix)
	}

	// Check if prefix contains date in YYYY-MM-DD format
	expectedDate := time.Now().Format("2006-01-02")
	if !strings.Contains(prefix, expectedDate) {
		t.Errorf("Expected prefix to contain date %s, got %s", expectedDate, prefix)
	}
}

func TestClose(t *testing.T) {
	log := NewLogger()

	// Close should not panic
	log.Close()
}

func TestShowCallerConfiguration(t *testing.T) {
	// Test with ShowCaller = true (default)
	configWithCaller := LoggerConfig{
		OutputMode: OutputTerminal,
		LogLevel:   LevelInfo,
		LogDir:     "test_logs",
		ShowCaller: true,
	}

	logWithCaller := NewLoggerWithConfig(configWithCaller)
	defer logWithCaller.Close()

	if !logWithCaller.showCaller {
		t.Error("Expected showCaller to be true")
	}

	// Test with ShowCaller = false
	configWithoutCaller := LoggerConfig{
		OutputMode: OutputTerminal,
		LogLevel:   LevelInfo,
		LogDir:     "test_logs",
		ShowCaller: false,
	}

	logWithoutCaller := NewLoggerWithConfig(configWithoutCaller)
	defer logWithoutCaller.Close()

	if logWithoutCaller.showCaller {
		t.Error("Expected showCaller to be false")
	}

	// Test default behavior (should be true)
	defaultLog := NewLogger()
	defer defaultLog.Close()

	if !defaultLog.showCaller {
		t.Error("Expected default showCaller to be true")
	}
}

func TestWithContextPreservesShowCaller(t *testing.T) {
	// Test that WithContext preserves the showCaller configuration
	config := LoggerConfig{
		OutputMode: OutputTerminal,
		LogLevel:   LevelInfo,
		LogDir:     "test_logs",
		ShowCaller: false,
	}

	log := NewLoggerWithConfig(config)
	defer log.Close()

	ctx := WithRequestID(context.Background(), "test-request")
	contextLogger := log.WithContext(ctx)

	if contextLogger.showCaller != log.showCaller {
		t.Errorf("Expected contextLogger.showCaller (%v) to match log.showCaller (%v)",
			contextLogger.showCaller, log.showCaller)
	}
}

// Benchmark tests
func BenchmarkSimpleLogging(b *testing.B) {
	log := NewLogger()
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("Benchmark message").Send()
	}
}

func BenchmarkLoggingWithData(b *testing.B) {
	log := NewLogger()
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("Benchmark message").
			Data("key1", "value1").
			Data("key2", 123).
			Data("key3", true).
			Send()
	}
}

func BenchmarkLoggingWithContext(b *testing.B) {
	log := NewLogger()
	defer log.Close()
	ctx := WithRequestID(context.Background(), "benchmark-request")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.WithContext(ctx).
			Info("Benchmark message").
			Data("iteration", i).
			Send()
	}
}

func BenchmarkMethodChaining(b *testing.B) {
	log := NewLogger()
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("Benchmark message").
			Data("key1", "value1").
			Data("key2", 123).
			Data("key3", true).
			Data("key4", 3.14).
			Data("key5", "another value").
			Send()
	}
}
