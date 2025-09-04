package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/risoftinc/gologger"
)

// HTTPExample demonstrates how to use the logger with HTTP requests
func HTTPExample() {
	// Create logger
	log := gologger.NewLogger()
	defer log.Close()

	// Setup HTTP handler
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handleUserRequest(w, r, log)
	})

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		handleHealthCheck(w, r, log)
	})

	fmt.Println("Server starting on :8080")
	fmt.Println("Try: curl http://localhost:8080/api/users")
	fmt.Println("Try: curl http://localhost:8080/api/health")

	// Start server (in real app, you'd use log.Fatal)
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUserRequest(w http.ResponseWriter, r *http.Request, log gologger.Logger) {
	// Generate request ID
	requestID := generateRequestID()

	// Add request ID to context
	ctx := gologger.WithRequestID(r.Context(), requestID)

	// Log request start
	log.WithContext(ctx).
		Info("HTTP request started").
		Data("method", r.Method).
		Data("path", r.URL.Path).
		Data("user_agent", r.UserAgent()).
		Data("remote_addr", r.RemoteAddr).
		Send()

	// Simulate processing
	user, err := processUserRequest(ctx, log, r)
	if err != nil {
		log.WithContext(ctx).
			Error("Failed to process user request").
			ErrorData(err).
			Send()
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log success
	log.WithContext(ctx).
		Info("HTTP request completed").
		Data("status", http.StatusOK).
		Data("user_id", user.ID).
		Data("processing_time", "150ms").
		Send()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id": %d, "name": "%s", "request_id": "%s"}`, user.ID, user.Name, requestID)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request, log gologger.Logger) {
	ctx := gologger.WithRequestID(r.Context(), "health-"+generateRequestID())

	log.WithContext(ctx).Info("Health check requested").Send()

	// Simulate health check
	status := checkSystemHealth(ctx, log)

	if status.Healthy {
		log.WithContext(ctx).
			Info("Health check passed").
			Data("response_time", status.ResponseTime).
			Send()
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "healthy", "timestamp": "`+time.Now().Format(time.RFC3339)+`"}`)
	} else {
		log.WithContext(ctx).
			Error("Health check failed").
			Data("error", status.Error).
			Send()
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"status": "unhealthy", "error": "`+status.Error+`"}`)
	}
}

func processUserRequest(ctx context.Context, log gologger.Logger, r *http.Request) (*User, error) {
	log.WithContext(ctx).
		Info("Processing user request").
		Data("step", "validation").
		Send()

	// Simulate validation
	if err := validateRequest(ctx, log, r); err != nil {
		return nil, err
	}

	log.WithContext(ctx).
		Info("Request validation passed").
		Data("step", "database").
		Send()

	// Simulate database operation
	user, err := saveUserToDatabase(ctx, log)
	if err != nil {
		return nil, err
	}

	log.WithContext(ctx).
		Info("User saved successfully").
		Data("user_id", user.ID).
		Send()

	return user, nil
}

func validateRequest(ctx context.Context, log gologger.Logger, r *http.Request) error {
	log.WithContext(ctx).
		Debug("Validating request").
		Data("component", "validator").
		Send()

	// Simulate validation logic
	if r.Method != "POST" {
		log.WithContext(ctx).
			Warn("Invalid HTTP method").
			Data("expected", "POST").
			Data("received", r.Method).
			Send()
		return fmt.Errorf("method not allowed")
	}

	log.WithContext(ctx).
		Debug("Request validation completed").
		Data("result", "success").
		Send()
	return nil
}

func saveUserToDatabase(ctx context.Context, log gologger.Logger) (*User, error) {
	log.WithContext(ctx).
		Info("Connecting to database").
		Data("component", "database").
		Send()

	// Simulate database connection
	time.Sleep(50 * time.Millisecond)

	log.WithContext(ctx).
		Debug("Executing INSERT query").
		Data("table", "users").
		Send()

	// Simulate database operation
	time.Sleep(100 * time.Millisecond)

	// Simulate potential database error
	if time.Now().UnixNano()%10 == 0 { // 10% chance of error
		log.WithContext(ctx).
			Error("Database operation failed").
			Data("error", "connection timeout").
			Send()
		return nil, fmt.Errorf("database connection timeout")
	}

	user := &User{
		ID:   12345,
		Name: "John Doe",
	}

	log.WithContext(ctx).
		Info("Database operation completed").
		Data("rows_affected", 1).
		Send()

	return user, nil
}

func checkSystemHealth(ctx context.Context, log gologger.Logger) *HealthStatus {
	log.WithContext(ctx).
		Debug("Checking system health").
		Data("component", "health_checker").
		Send()

	// Simulate health checks
	start := time.Now()

	// Check database
	if err := checkDatabase(ctx, log); err != nil {
		return &HealthStatus{
			Healthy:      false,
			Error:        "database unavailable",
			ResponseTime: time.Since(start).String(),
		}
	}

	// Check external services
	if err := checkExternalServices(ctx, log); err != nil {
		return &HealthStatus{
			Healthy:      false,
			Error:        "external services unavailable",
			ResponseTime: time.Since(start).String(),
		}
	}

	log.WithContext(ctx).
		Debug("All health checks passed").
		Data("duration", time.Since(start)).
		Send()

	return &HealthStatus{
		Healthy:      true,
		ResponseTime: time.Since(start).String(),
	}
}

func checkDatabase(ctx context.Context, log gologger.Logger) error {
	log.WithContext(ctx).
		Debug("Checking database connection").
		Data("service", "database").
		Send()
	time.Sleep(10 * time.Millisecond)
	return nil
}

func checkExternalServices(ctx context.Context, log gologger.Logger) error {
	log.WithContext(ctx).
		Debug("Checking external services").
		Data("service", "payment_gateway").
		Send()
	time.Sleep(20 * time.Millisecond)
	return nil
}

func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

// Data structures
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type HealthStatus struct {
	Healthy      bool   `json:"healthy"`
	Error        string `json:"error,omitempty"`
	ResponseTime string `json:"response_time"`
}
