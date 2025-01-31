package main

import "testing"

// Potential Client Test Cases

// Case 1: Version on the server is same as client
// Expected Result: No attempt to download and client message app is up-to-date
func TestVersionMatch(t *testing.T) {
	// Test logic
}

// Case 2: Bad Binary fails to update
// Expected Result: Program reverts to previous working version, and retries 3 times
func TestBadBinary(t *testing.T) {
	// Test logic
}

// Integration Tests

// Case 1: Main client update
// Expected Result: Program initial version is replaced with pred-defined update version
func TestUpdateFlow(t *testing.T) {
	// Test logic
}
