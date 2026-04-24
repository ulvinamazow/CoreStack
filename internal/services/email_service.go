package services

import (
	"fmt"
	"net/smtp"

	"github.com/ulvinamazow/CoreStack/internal/config"
)

func SendVerificationEmail(toEmail, toName, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", config.App.AppURL, token)

	subject := "Verify your CoreStack account"
	body := fmt.Sprintf(`
Hello %s,

Thank you for registering at CoreStack!

Please verify your email address by clicking the link below:
%s

This link will expire in %d hours.

If you did not create an account, please ignore this email.

Best regards,
The CoreStack Team
`, toName, verifyURL, config.App.VerificationTokenExpiryHours)

	return sendEmail(toEmail, subject, body)
}

func SendPasswordResetEmail(toEmail, toName, resetURL string) error {
	subject := "Reset your CoreStack password"
	body := fmt.Sprintf(`
Hello %s,

You requested a password reset for your CoreStack account.

Click the link below to reset your password:
%s

If you did not request this, please ignore this email.

Best regards,
The CoreStack Team
`, toName, resetURL)

	return sendEmail(toEmail, subject, body)
}

func sendEmail(to, subject, body string) error {
	cfg := config.App
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		cfg.SMTPUser, to, subject, body)

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, cfg.SMTPUser, []string{to}, []byte(msg))
}
