package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("🧪 Running Logger Package Tests...")
	fmt.Println(strings.Repeat("=", 50))

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	fmt.Println("📋 Running unit tests...")
	cmd := exec.Command("go", "test", "-v", "-cover", "./...")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		fmt.Printf("❌ Tests failed:\n%s\n", string(output))
		os.Exit(1)
	}

	fmt.Printf("✅ Tests passed:\n%s\n", string(output))

	// Run benchmarks
	fmt.Println("\n📊 Running benchmarks...")
	cmd = exec.Command("go", "test", "-bench=.", "-benchmem", "./...")
	cmd.Dir = dir
	output, err = cmd.CombinedOutput()
	
	if err != nil {
		fmt.Printf("❌ Benchmarks failed:\n%s\n", string(output))
		os.Exit(1)
	}

	fmt.Printf("✅ Benchmarks completed:\n%s\n", string(output))

	// Run race detection
	fmt.Println("\n🔍 Running race detection...")
	cmd = exec.Command("go", "test", "-race", "./...")
	cmd.Dir = dir
	output, err = cmd.CombinedOutput()
	
	if err != nil {
		fmt.Printf("❌ Race detection failed:\n%s\n", string(output))
		os.Exit(1)
	}

	fmt.Printf("✅ No race conditions detected:\n%s\n", string(output))

	// Clean up test files
	fmt.Println("\n🧹 Cleaning up test files...")
	cleanupTestFiles(dir)

	fmt.Println("\n🎉 All tests completed successfully!")
}

func cleanupTestFiles(dir string) {
	testDirs := []string{"test_logs", "logs"}
	
	for _, testDir := range testDirs {
		fullPath := filepath.Join(dir, testDir)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			err := os.RemoveAll(fullPath)
			if err != nil {
				fmt.Printf("⚠️  Warning: Could not remove %s: %v\n", fullPath, err)
			} else {
				fmt.Printf("✅ Removed %s\n", fullPath)
			}
		}
	}
}
