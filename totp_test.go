package totp

import (
	"testing"
	"time"
)

func TestGenerateSecret(t *testing.T) {
	// Test generating a secret
	secret, err := GenerateSecret()
	if err != nil {
		t.Fatalf("Failed to generate secret: %v", err)
	}

	// Check secret length (base32 encoded, should be around 16 characters)
	if len(secret) < 10 || len(secret) > 20 {
		t.Errorf("Generated secret length is invalid: %d", len(secret))
	}
}

func TestTOTP(t *testing.T) {
	// Test scenarios
	testCases := []struct {
		name     string
		duration int
		wantErr  bool
	}{
		{"Standard 30-second interval", 30, false},
		{"Custom 60-second interval", 60, false},
		{"Invalid duration", -1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate a secret
			secret, err := GenerateSecret()
			if err != nil {
				t.Fatalf("Failed to generate secret: %v", err)
			}

			// Generate TOTP
			code, err := TOTP(secret, tc.duration)
			if tc.wantErr && err == nil {
				t.Error("Expected an error, but got none")
			}

			if err != nil && !tc.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check code length
			if len(code) != 6 {
				t.Errorf("Invalid TOTP code length: %d", len(code))
			}
		})
	}
}

// timeProvider is an interface to allow easier time mocking
type timeProvider interface {
	Now() time.Time
}

// realTimeProvider uses the actual system time
type realTimeProvider struct{}

func (r realTimeProvider) Now() time.Time {
	return time.Now()
}

// mockTimeProvider allows setting a fixed time for testing
type mockTimeProvider struct {
	fixedTime time.Time
}

func (m mockTimeProvider) Now() time.Time {
	return m.fixedTime
}

// Global variable to hold current time provider
var currentTimeProvider timeProvider = realTimeProvider{}

// SetTimeProvider allows setting a custom time provider (useful for testing)
func SetTimeProvider(provider timeProvider) {
	currentTimeProvider = provider
}

// Modified functions to use the time provider
func nowFunc() time.Time {
	return currentTimeProvider.Now()
}


// Updated test file
func TestValidate(t *testing.T) {
	// Reset time provider after the test
	defer SetTimeProvider(realTimeProvider{})

	// Test validation scenarios
	testCases := []struct {
		name     string
		duration int
		timeDiff time.Duration
		expected bool
	}{
		{"Current time", 30, 0, true},
		{"30 seconds before", 30, -30 * time.Second, true},
		{"30 seconds after", 30, 30 * time.Second, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fixed time and set it as the current time provider
			fixedTime := time.Now().Add(tc.timeDiff)
			SetTimeProvider(mockTimeProvider{fixedTime})

			// Generate a secret
			secret, err := GenerateSecret()
			if err != nil {
				t.Fatalf("Failed to generate secret: %v", err)
			}

			// Generate TOTP code
			code, err := TOTP(secret, tc.duration)
			if err != nil {
				t.Fatalf("Failed to generate TOTP: %v", err)
			}

			// Validate the code
			isValid := Validate(secret, tc.duration, code)
			if isValid != tc.expected {
				t.Errorf("Validation result unexpected. Got %v, want %v", isValid, tc.expected)
			}
		})
	}
}

func TestInvalidSecret(t *testing.T) {
	// Test with invalid secrets
	invalidSecrets := []string{
		"INVALID_SECRET",
		"12345",
	}

	for _, secret := range invalidSecrets {
		t.Run("Invalid Secret: "+secret, func(t *testing.T) {
			// Try to generate TOTP with invalid secret
			_, err := TOTP(secret, 30)
			if err == nil {
				t.Error("Expected an error with invalid secret, but got none")
			}

			// Try to validate with invalid secret
			isValid := Validate(secret, 30, "123456")
			if isValid {
				t.Error("Validation should fail with invalid secret")
			}
		})
	}
}

// Example of resetting to real time provider
func resetTimeProvider() {
	SetTimeProvider(realTimeProvider{})
}

func BenchmarkGenerateSecret(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateSecret()
	}
}

func BenchmarkTOTP(b *testing.B) {
	secret, _ := GenerateSecret()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = TOTP(secret, 30)
	}
}

func BenchmarkValidate(b *testing.B) {
	secret, _ := GenerateSecret()
	code, _ := TOTP(secret, 30)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Validate(secret, 30, code)
	}
}

// Example of how to use the package
func ExampleTOTP() {
	// Generate a secret
	secret, err := GenerateSecret()
	if err != nil {
		panic(err)
	}

	// Generate a TOTP code
	code, err := TOTP(secret, 30)
	if err != nil {
		panic(err)
	}

	// Validate the code
	isValid := Validate(secret, 30, code)
	println("Code is valid:", isValid)
}
