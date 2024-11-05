package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

)

// GenerateSecret generates a random TOTP secret key in base32 encoding.
func GenerateSecret() (string, error) {
	
    secretLength := 10

    // Define the length of the random byte sequence (usually between 10 to 20 bytes for TOTP)
	randomBytes := make([]byte, secretLength)

    
	// Generate random bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("error generating random bytes: %v", err)
	}

	// Encode the random bytes to base32
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	return secret, nil
}

// TOTP generates a 6-digit TOTP code using the given base32-encoded secret and a time step of 30 seconds.
func TOTP(secret string, duration int) (string, error) {
	// Decode the base32-encoded secret
	secret = strings.ToUpper(secret) // TOTP secrets are usually upper-case

	// Calculate the time counter, which is the number of 30-second intervals since Unix epoch
	timestamp := int64(time.Now().Unix())

    return generateTOTP(secret, timestamp, duration)
}

// ValidateTOTP checks if the provided code matches the generated TOTP code for the given secret
func Validate(secret string, duration int, code string) bool {
    if duration < 1 {
        duration = 30
    }
	// Allow for a +/- 30-second time window to account for clock drift
	currentTimestamp := time.Now().Unix()
	for i := -1; i <= 1; i++ {
		// Generate the TOTP for the current time window
        timestamp := int64(currentTimestamp + int64(i * duration))
		generatedCode, err := generateTOTP(secret, timestamp, duration)
		if err != nil {
			fmt.Println("Error generating TOTP:", err)
			return false
		}

		// Check if the generated TOTP matches the provided code
		if generatedCode == code {
			return true
		}
	}

	return false
}

// generateTOTP generates a TOTP code for the current timestamp
func generateTOTP(secret string, timestamp int64, duration int) (string, error) {
	// Decode the base32-encoded secret
	secret = strings.ToUpper(secret) // TOTP secrets are usually upper-case
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("error decoding secret: %v", err)
	}

	// Calculate the time counter, based on the timestamp
	interval := uint64(timestamp / int64(duration))

	// Convert the time counter to a byte array (big-endian)
	var counterBytes [8]byte
	binary.BigEndian.PutUint64(counterBytes[:], interval)

	// Create an HMAC-SHA1 hash using the counter as the message and the secret as the key
	hmacHash := hmac.New(sha1.New, key)
	hmacHash.Write(counterBytes[:])
	hash := hmacHash.Sum(nil)

	// Perform dynamic truncation to get a 4-byte string (as per TOTP spec)
	offset := hash[len(hash)-1] & 0x0F
	code := (int(hash[offset]&0x7F) << 24) |
		(int(hash[offset+1]&0xFF) << 16) |
		(int(hash[offset+2]&0xFF) << 8) |
		(int(hash[offset+3] & 0xFF))

	// Get the last 6 digits of the code as the OTP
	otp := code % 1000000

	// Return the OTP as a zero-padded 6-digit string
	return fmt.Sprintf("%06d", otp), nil
}

