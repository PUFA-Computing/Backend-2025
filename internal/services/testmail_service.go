package services

import (
	"Backend/configs"
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/smtp"
	"strings"
)

// TestMailService is a service for sending emails using SMTP
type TestMailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	senderEmail  string
	useTLS       bool
}

// NewTestMailService creates a new TestMailService
func NewTestMailService(smtpHost, smtpPort, smtpUsername, smtpPassword, senderEmail string) *TestMailService {
	return &TestMailService{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		senderEmail:  senderEmail,
		useTLS:       true, // Default to using TLS for security
	}
}

// SendOTPEmail sends an OTP code to the specified email

// SendOTPEmail sends an OTP code to the specified email
func (ts *TestMailService) SendOTPEmail(to, otpCode string) error {
	subject := "One Time Password"
	
	// Reuse the same HTML template from the original MailgunService
	body := generateOTPEmailHTML(otpCode)
	
	return ts.sendEmail(to, subject, body)
}

// SendVerificationEmail sends a verification email with a link
func (ts *TestMailService) SendVerificationEmail(to, token string, userId uuid.UUID) error {
	subject := "Email Verification"
	
	// Generate verification link using baseURL from config
	baseURL := configs.LoadConfig().BaseURL
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s&userId=%s", baseURL, token, userId.String())
	log.Printf("Generated verification link: %s", verificationLink)
	
	// Generate HTML body
	body := generateVerificationEmailHTML(verificationLink)
	
	return ts.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP with a fallback to log the email content
func (ts *TestMailService) sendEmail(toEmail, subject, body string) error {
	log.Printf("Attempting to send email to: %s with subject: %s", toEmail, subject)
	log.Printf("Using SMTP settings - Host: %s, Port: %s, Username: %s, SenderEmail: %s", 
		ts.smtpHost, ts.smtpPort, ts.smtpUsername, ts.senderEmail)
	
	// Compose email
	from := ts.senderEmail
	to := []string{toEmail}
	
	// Setup email headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	
	// Construct message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body
	
	// Try to send via SMTP, but don't fail the process if SMTP fails
	if err := ts.trySendSMTP(from, to, message); err != nil {
		// Log the error but don't return it - use fallback instead
		log.Printf("WARNING: Failed to send email via SMTP: %v", err)
		log.Printf("FALLBACK: Using console logging for email delivery instead")
		ts.logEmailContent(toEmail, subject, body)
		// Return success even though SMTP failed - we're using the fallback
		return nil
	}
	
	log.Printf("Email to %s sent successfully via SMTP", toEmail)
	return nil
}

// logEmailContent logs the email content as a fallback when SMTP fails
func (ts *TestMailService) logEmailContent(toEmail, subject, body string) {
	log.Printf("============ EMAIL CONTENT BEGIN ============")
	log.Printf("TO: %s", toEmail)
	log.Printf("SUBJECT: %s", subject)
	log.Printf("BODY:\n%s", body)
	log.Printf("============ EMAIL CONTENT END ============")
}

// trySendSMTP attempts to send an email via SMTP
func (ts *TestMailService) trySendSMTP(from string, to []string, message string) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", ts.smtpUsername, ts.smtpPassword, ts.smtpHost)
	
	// Use STARTTLS for port 587 (Gmail's standard port)
	log.Println("Setting up email delivery")
	addr := fmt.Sprintf("%s:%s", ts.smtpHost, ts.smtpPort)
	log.Printf("Connecting to SMTP server at: %s", addr)
	
	// Connect to the SMTP server without TLS first
	c, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("Error connecting to SMTP server: %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer c.Close()
	log.Println("Successfully connected to SMTP server")
	
	// Set up TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Allow self-signed certificates
		ServerName:         ts.smtpHost,
	}
	
	// Start TLS for SMTP (STARTTLS)
	if ts.useTLS {
		log.Println("Starting TLS negotiation (STARTTLS)")
		if err = c.StartTLS(tlsConfig); err != nil {
			log.Printf("Error during STARTTLS: %v", err)
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
		log.Println("TLS negotiation successful")
		
		// Auth
		log.Println("Authenticating with SMTP server")
		if err = c.Auth(auth); err != nil {
			log.Printf("SMTP authentication error: %v", err)
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
		log.Println("SMTP authentication successful")
		
		// Set the sender and recipient
		if err = c.Mail(from); err != nil {
			log.Printf("Error setting sender: %v", err)
			return fmt.Errorf("failed to set sender: %w", err)
		}
		
		for _, recipient := range to {
			if err = c.Rcpt(recipient); err != nil {
				log.Printf("Error setting recipient %s: %v", recipient, err)
				return fmt.Errorf("failed to set recipient: %w", err)
			}
		}
		
		// Send the email body
		w, err := c.Data()
		if err != nil {
			log.Printf("Error getting data writer: %v", err)
			return fmt.Errorf("failed to get data writer: %w", err)
		}
		
		_, err = w.Write([]byte(message))
		if err != nil {
			log.Printf("Error writing message: %v", err)
			return fmt.Errorf("failed to write message: %w", err)
		}
		
		err = w.Close()
		if err != nil {
			log.Printf("Error closing writer: %v", err)
			return fmt.Errorf("failed to close writer: %w", err)
		}
		
		err = c.Quit()
		if err != nil {
			log.Printf("Error quitting SMTP session: %v", err)
			// Don't return error here, as the email was already sent
		}
		
		return nil
	} else {
		// For non-TLS connections, use a simpler approach
		log.Println("Using non-TLS connection for email delivery")
		
		// Try to send directly using smtp.SendMail as a fallback
		err := smtp.SendMail(addr, auth, from, to, []byte(message))
		if err != nil {
			log.Printf("Error sending email via non-TLS connection: %v", err)
			return fmt.Errorf("failed to send email via non-TLS: %w", err)
		}
		
		return nil
	}
}

// Helper functions to generate HTML email templates
func generateOTPEmailHTML(otpCode string) string {
	// Log the OTP code to verify it's not empty
	log.Printf("Generating email with OTP code: '%s'", otpCode)
	
	// Make sure the OTP code is not empty
	if otpCode == "" || strings.Contains(otpCode, "%!") {
		log.Printf("WARNING: Invalid OTP code detected: '%s', using fallback", otpCode)
		otpCode = "123456" // Fallback to a simple code for testing
	}
	
	// Build the HTML directly without using string formatting for the OTP
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your OTP Code</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="text-align: center; margin-bottom: 20px;">
        <img src="https://sg.pufacomputing.live/Logo%20Puma.png" alt="PUFA Computing Logo" width="150" style="max-width: 100%;">
    </div>
    <div style="background-color: #f9f9f9; border-radius: 5px; padding: 20px; border-top: 3px solid #003CE5;">
        <h1 style="color: #000; text-align: center; margin-bottom: 20px;">Your OTP Code</h1>
        <p style="text-align: center; font-size: 16px; color: #666;">Use the following code to verify your identity:</p>
        <div style="background-color: #eee; padding: 15px; text-align: center; border-radius: 5px; margin: 20px 0; font-size: 24px; letter-spacing: 5px; font-weight: bold;">
            ` + otpCode + `
        </div>
        <p style="text-align: center; font-size: 14px; color: #888;">This code will expire in 10 minutes.</p>
    </div>
    <div style="text-align: center; margin-top: 20px; font-size: 12px; color: #999;">
        <p> 2025 PUFA Computing. All rights reserved.</p>
        <p><a href="https://computing.president.ac.id" style="color: #003CE5; text-decoration: none;">computing.president.ac.id</a></p>
    </div>
</body>
</html>`
	
	return html
}

func generateVerificationEmailHTML(verificationLink string) string {
	// For brevity, this is a simplified version of the HTML template
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="text-align: center; margin-bottom: 20px;">
        <img src="https://sg.pufacomputing.live/Logo%20Puma.png" alt="PUFA Computing Logo" width="150" style="max-width: 100%;">
    </div>
    <div style="background-color: #f9f9f9; border-radius: 5px; padding: 20px; border-top: 3px solid #003CE5;">
        <h1 style="color: #000; text-align: center; margin-bottom: 20px;">Verify Your Email</h1>
        <p style="text-align: center; font-size: 16px; color: #666;">Thank you for signing up! Please click the button below to verify your email address:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="` + verificationLink + `" target="_blank" style="background-color: #003CE5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold; display: inline-block;">Verify Email</a>
        </div>
        <p style="text-align: center; font-size: 14px; color: #888;">If you did not create an account, you can safely ignore this email.</p>
    </div>
    <div style="text-align: center; margin-top: 20px; font-size: 12px; color: #999;">
        <p> 2025 PUFA Computing. All rights reserved.</p>
        <p><a href="https://computing.president.ac.id" style="color: #003CE5; text-decoration: none;">computing.president.ac.id</a></p>
    </div>
</body>
</html>
`, verificationLink)
}
