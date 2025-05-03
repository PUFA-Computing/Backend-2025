package services

import (
	"Backend/configs"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridService is a service for sending emails using SendGrid
type SendGridService struct {
	apiKey      string
	senderEmail string
	senderName  string
}

// NewSendGridService creates a new SendGridService
func NewSendGridService(apiKey, senderEmail, senderName string) *SendGridService {
	return &SendGridService{
		apiKey:      apiKey,
		senderEmail: senderEmail,
		senderName:  senderName,
	}
}

// SendOTPEmail sends an OTP code to the specified email
func (sg *SendGridService) SendOTPEmail(to, otpCode string) error {
	subject := "One Time Password"

	// Use our own HTML template with proper formatting
	body := sg.generateOTPEmailHTML(otpCode)

	return sg.sendEmail(to, subject, body)
}

// SendVerificationEmail sends a verification email with a link
func (sg *SendGridService) SendVerificationEmail(to, token string, userId uuid.UUID) error {
	subject := "Email Verification"

	// Generate verification link using baseURL from config
	baseURL := configs.LoadConfig().BaseURL
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s&userId=%s", baseURL, token, userId.String())
	log.Printf("Generated verification link: %s", verificationLink)

	// Generate HTML body
	body := generateVerificationEmailHTML(verificationLink)

	return sg.sendEmail(to, subject, body)
}

// sendEmail sends an email using SendGrid
func (sg *SendGridService) sendEmail(toEmail, subject, htmlContent string) error {
	log.Printf("Attempting to send email to: %s with subject: %s", toEmail, subject)
	log.Printf("Using SendGrid with sender: %s <%s>", sg.senderName, sg.senderEmail)

	// Add debug logging for API key (partial, for security)
	if sg.apiKey == "" {
		log.Printf("ERROR: SendGrid API key is empty!")
	} else {
		keyLength := len(sg.apiKey)
		log.Printf("SendGrid API key found (length: %d, first 4 chars: %s...)", keyLength, sg.apiKey[:4])
	}

	from := mail.NewEmail(sg.senderName, sg.senderEmail)
	to := mail.NewEmail("", toEmail)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(sg.apiKey)

	log.Printf("Sending email via SendGrid API...")
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Error sending email via SendGrid: %v", err)
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	log.Printf("SendGrid response - StatusCode: %d, Body: %s", response.StatusCode, response.Body)

	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid API error - StatusCode: %d, Body: %s", response.StatusCode, response.Body)
	}

	log.Printf("Email to %s sent successfully via SendGrid", toEmail)
	return nil
}

// generateOTPEmailHTML creates HTML content for OTP emails with proper formatting
func (sg *SendGridService) generateOTPEmailHTML(otpCode string) string {
	// For brevity, this is a simplified version of the HTML template
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your OTP Code</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="text-align: center; margin-bottom: 20px;">
        <img src="https://sg.pufacomputing.live/Logo%%20Puma.png" alt="PUFA Computing Logo" width="150" style="max-width: 100%%;">
    </div>
    <div style="background-color: #f9f9f9; border-radius: 5px; padding: 20px; border-top: 3px solid #003CE5;">
        <h1 style="color: #000; text-align: center; margin-bottom: 20px;">Your OTP Code</h1>
        <p style="text-align: center; font-size: 16px; color: #666;">Use the following code to verify your identity:</p>
        <div style="background-color: #eee; padding: 15px; text-align: center; border-radius: 5px; margin: 20px 0; font-size: 24px; letter-spacing: 5px; font-weight: bold;">
            %s
        </div>
        <p style="text-align: center; font-size: 14px; color: #888;">This code will expire in 10 minutes.</p>
    </div>
    <div style="text-align: center; margin-top: 20px; font-size: 12px; color: #999;">
        <p> 2025 PUFA Computing. All rights reserved.</p>
        <p><a href="https://compsci.president.ac.id" style="color: #003CE5; text-decoration: none;">compsci.president.ac.id</a></p>
    </div>
</body>
</html>
`, otpCode)
}
