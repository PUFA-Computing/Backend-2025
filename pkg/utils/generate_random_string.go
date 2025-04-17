package utils

import (
	cryptoRand "crypto/rand"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func init() {
	// Seed the math/rand package with the current time
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

const alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomOTPCode generates a random 6-digit OTP code
func GenerateRandomOTPCode() string {
	// Use a more reliable method to generate numeric OTP
	const digits = "0123456789"
	const otpLength = 6
	
	// Create a buffer to hold the result
	result := make([]byte, otpLength)
	
	// Fill the buffer with random digits
	for i := 0; i < otpLength; i++ {
		// Get a random byte
		byte := make([]byte, 1)
		_, err := cryptoRand.Read(byte)
		if err != nil {
			// Fallback to a simpler method if crypto/rand fails
			result[i] = digits[rand.Intn(10)]
			continue
		}
		
		// Convert the random byte to a digit (0-9)
		result[i] = digits[int(byte[0])%10]
	}
	
	// Convert to string and log for debugging
	otpCode := string(result)
	log.Printf("Generated OTP code: %s", otpCode)
	
	// Ensure we don't return an empty string
	if otpCode == "" {
		// Fallback to a simple timestamp-based code if all else fails
		otpCode = fmt.Sprintf("%06d", time.Now().Unix()%1000000)
		log.Printf("Used fallback OTP generation: %s", otpCode)
	}
	
	return otpCode
}

// GenerateRandomTokenOtp is a secret that is used to generate a random OTP
func GenerateRandomTokenOtp() string {
	// Generate cryptographically secure random bytes
	token := make([]byte, 32)
	_, err := cryptoRand.Read(token)
	if err != nil {
		// Log the error
		log.Printf("Error generating random token: %v", err)
		// Fallback to a less secure but functional alternative
		for i := 0; i < 32; i++ {
			token[i] = byte(rand.Intn(256))
		}
	}

	// Filter out non-alphanumeric characters using a mask
	var mask byte = 0b111111 // Mask to keep only alphanumeric characters (lower and uppercase letters, digits)
	for i := range token {
		token[i] &= mask
	}

	// Select random characters from the alphanumeric charset
	for i, b := range token {
		token[i] = alphanumericCharset[b%byte(len(alphanumericCharset))]

		// Ensure that the generated token is 32 characters long
		if i == 31 {
			break
		}
	}

	return string(token)
}
