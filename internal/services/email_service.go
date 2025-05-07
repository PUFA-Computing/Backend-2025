package services

import (
	"github.com/google/uuid"
)

// EmailService is an interface for email services
type EmailService interface {
	// SendOTPEmail sends an OTP code to the specified email
	SendOTPEmail(to, otpCode string) error
	
	// SendVerificationEmail sends a verification email with a link
	SendVerificationEmail(to, token string, userId uuid.UUID) error
}
