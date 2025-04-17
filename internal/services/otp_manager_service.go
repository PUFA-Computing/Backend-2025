package services

import (
	"Backend/pkg/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type OTPManager struct {
}

func NewOTPManager() *OTPManager {
	return &OTPManager{}
}

// In-memory fallback storage for OTPs when Redis is unavailable
var memoryOTPStore = make(map[string]map[string]string)

func (om *OTPManager) GenerateOTP(userID uuid.UUID, token string, expiration time.Duration) (string, error) {
	otpCode := utils.GenerateRandomOTPCode()
	expiresAt := time.Now().Add(expiration)

	// Try to store in Redis first
	ctx := context.Background()
	otpKey := "otp:" + userID.String()

	// Try Redis operation with proper error handling
	redisErr := om.storeOTPInRedis(ctx, otpKey, token, otpCode, expiresAt, expiration)
	if redisErr != nil {
		// If Redis fails, log the error and use in-memory fallback
		log.Printf("Redis error in GenerateOTP: %v. Using in-memory fallback.", redisErr)
		
		// Use in-memory fallback
		om.storeOTPInMemory(otpKey, token, otpCode, expiresAt)
	}

	log.Printf("Generated OTP code for user %s: %s", userID.String(), otpCode)
	return otpCode, nil
}

// storeOTPInRedis attempts to store OTP data in Redis
func (om *OTPManager) storeOTPInRedis(ctx context.Context, otpKey, token, otpCode string, expiresAt time.Time, expiration time.Duration) error {
	// First check if Redis is available by pinging
	_, pingErr := utils.Rdb.Ping(ctx).Result()
	if pingErr != nil {
		return fmt.Errorf("Redis not available: %w", pingErr)
	}

	// Try to store the OTP data
	try := func() error {
		err := utils.Rdb.HSet(ctx, otpKey, 
			"TokenOTP", token,
			"OTPCode", otpCode,
			"ExpiresAt", expiresAt.Format(time.RFC3339),
		).Err()
		if err != nil {
			return err
		}

		err = utils.Rdb.Expire(ctx, otpKey, expiration).Err()
		if err != nil {
			return err
		}

		return nil
	}

	// Try the operation
	err := try()
	if err != nil {
		// If it fails, try an alternative approach
		log.Printf("First Redis attempt failed: %v. Trying alternative approach.", err)
		
		// Try setting each field individually
		err = utils.Rdb.HSet(ctx, otpKey, "TokenOTP", token).Err()
		if err == nil {
			err = utils.Rdb.HSet(ctx, otpKey, "OTPCode", otpCode).Err()
		}
		if err == nil {
			err = utils.Rdb.HSet(ctx, otpKey, "ExpiresAt", expiresAt.Format(time.RFC3339)).Err()
		}
		if err == nil {
			err = utils.Rdb.Expire(ctx, otpKey, expiration).Err()
		}
	}

	return err
}

// storeOTPInMemory stores OTP data in memory as a fallback
func (om *OTPManager) storeOTPInMemory(otpKey, token, otpCode string, expiresAt time.Time) {
	memoryOTPStore[otpKey] = map[string]string{
		"TokenOTP":  token,
		"OTPCode":   otpCode,
		"ExpiresAt": expiresAt.Format(time.RFC3339),
	}

	// Set up a goroutine to clean up expired OTPs
	go func() {
		time.Sleep(time.Until(expiresAt))
		delete(memoryOTPStore, otpKey)
	}()
}

func (om *OTPManager) VerifyOTP(userID uuid.UUID, tokenOtp, otpCode string) bool {
	ctx := context.Background()
	otpKey := "otp:" + userID.String()

	// Try to verify using Redis first
	verified, redisErr := om.verifyOTPWithRedis(ctx, otpKey, tokenOtp, otpCode)
	if redisErr != nil {
		// If Redis fails, log the error and try in-memory fallback
		log.Printf("Redis error in VerifyOTP: %v. Trying in-memory fallback.", redisErr)
		return om.verifyOTPWithMemory(otpKey, tokenOtp, otpCode)
	}

	return verified
}

// verifyOTPWithRedis attempts to verify OTP using Redis
func (om *OTPManager) verifyOTPWithRedis(ctx context.Context, otpKey, tokenOtp, otpCode string) (bool, error) {
	// First check if Redis is available by pinging
	_, pingErr := utils.Rdb.Ping(ctx).Result()
	if pingErr != nil {
		return false, fmt.Errorf("Redis not available: %w", pingErr)
	}

	otpData, err := utils.Rdb.HGetAll(ctx, otpKey).Result()
	if err != nil {
		return false, fmt.Errorf("error retrieving OTP data: %w", err)
	}

	if len(otpData) == 0 {
		return false, nil // No error, but OTP not found
	}

	storedTokenOTP := otpData["TokenOTP"]
	storedOTPCode := otpData["OTPCode"]
	expiresAt, err := time.Parse(time.RFC3339, otpData["ExpiresAt"])
	if err != nil {
		return false, fmt.Errorf("error parsing OTP expiration time: %w", err)
	}

	if storedTokenOTP != tokenOtp {
		return false, nil
	}

	if time.Now().After(expiresAt) {
		return false, nil
	}

	return storedOTPCode == otpCode, nil
}

// verifyOTPWithMemory verifies OTP using the in-memory store
func (om *OTPManager) verifyOTPWithMemory(otpKey, tokenOtp, otpCode string) bool {
	// Check if we have this OTP in memory
	otpData, exists := memoryOTPStore[otpKey]
	if !exists {
		return false
	}

	storedTokenOTP := otpData["TokenOTP"]
	storedOTPCode := otpData["OTPCode"]
	expiresAt, err := time.Parse(time.RFC3339, otpData["ExpiresAt"])
	if err != nil {
		log.Printf("Error parsing in-memory OTP expiration time: %v", err)
		return false
	}

	if storedTokenOTP != tokenOtp {
		return false
	}

	if time.Now().After(expiresAt) {
		// Clean up expired OTP
		delete(memoryOTPStore, otpKey)
		return false
	}

	return storedOTPCode == otpCode
}
