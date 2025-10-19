package services

import (
	"fmt"
	"strings"
)

const (
	defaultWebsiteURL  = "https://parallax.com"
	defaultFacebookURL = "https://www.facebook.com/parallax"
)

type IEmailTemplates interface {
	GetNewUserEmailTemplate(baseURL string, verificationToken string) string
	GetLoginEmailTemplate(baseURL string, verificationToken string) string
	GetPasswordResetEmailTemplate(baseURL string, verificationToken string) string
}

type EmailTemplates struct {
}

func NewEmailTemplates() IEmailTemplates {
	return &EmailTemplates{}
}

func (e *EmailTemplates) GetNewUserEmailTemplate(baseURL string, verificationToken string) string {
	sanitizedBaseURL := e.normalizeBaseURL(baseURL)
	verifyURL := fmt.Sprintf("%s/api/auth/verify?token=%s", sanitizedBaseURL, verificationToken)
	if sanitizedBaseURL == "" {
		verifyURL = fmt.Sprintf("/api/auth/verify?token=%s", verificationToken)
	}

	mainContent := fmt.Sprintf(`
		<h1 style="margin:0 0 16px 0; font-size:28px; font-weight:700; color:#0f172a;">Welcome to Smarter Lynx!</h1>
		<p style="margin:0 0 24px 0; font-size:16px; line-height:1.6; color:#1f2937;">
			Hello and Welcome to Smarter Lynx! Please verify your account by using the secure button below.
		</p>
		<table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 auto 24px auto;">
			<tr>
				<td align="center" style="border-radius:999px; background-color:#2563eb;">
					<a href="%s" style="display:inline-block; padding:14px 32px; font-size:16px; font-weight:600; color:#ffffff; text-decoration:none; border-radius:999px;">Verify Account</a>
				</td>
			</tr>
		</table>
		<p style="margin:0 0 8px 0; font-size:14px; line-height:1.6; color:#475569;">
			Or copy and paste this link into your browser:
		</p>
		<p style="margin:0 0 24px 0; font-size:14px; line-height:1.6;">
			<a href="%s" style="color:#2563eb; word-break:break-all;">%s</a>
		</p>
		<p style="margin:0; font-size:14px; line-height:1.6; color:#475569;">
			We're excited to have you with us. If you didn't request this email, you can safely ignore it.
		</p>
	`, verifyURL, verifyURL, verifyURL)

	return e.wrapEmail(sanitizedBaseURL, mainContent)
}

func (e *EmailTemplates) GetLoginEmailTemplate(baseURL string, verificationToken string) string {
	sanitizedBaseURL := e.normalizeBaseURL(baseURL)
	loginURL := fmt.Sprintf("%s/api/auth/login-with-email?token=%s", sanitizedBaseURL, verificationToken)
	if sanitizedBaseURL == "" {
		loginURL = fmt.Sprintf("/api/auth/login-with-email?token=%s", verificationToken)
	}

	mainContent := fmt.Sprintf(`
		<h1 style="margin:0 0 16px 0; font-size:28px; font-weight:700; color:#0f172a;">Log in to your account</h1>
		<p style="margin:0 0 24px 0; font-size:16px; line-height:1.6; color:#1f2937;">
			Hello! You can login to your account by using the secure button below.
		</p>
		<table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 auto 24px auto;">
			<tr>
				<td align="center" style="border-radius:999px; background-color:#2563eb;">
					<a href="%s" style="display:inline-block; padding:14px 32px; font-size:16px; font-weight:600; color:#ffffff; text-decoration:none; border-radius:999px;">Log In Instantly</a>
				</td>
			</tr>
		</table>
		<p style="margin:0 0 8px 0; font-size:14px; line-height:1.6; color:#475569;">
			Or copy and paste this link into your browser:
		</p>
		<p style="margin:0 0 24px 0; font-size:14px; line-height:1.6;">
			<a href="%s" style="color:#2563eb; word-break:break-all;">%s</a>
		</p>
		<p style="margin:0; font-size:14px; line-height:1.6; color:#475569;">
			If you did not request this email, please ignore it.
		</p>
	`, loginURL, loginURL, loginURL)

	return e.wrapEmail(sanitizedBaseURL, mainContent)
}

func (e *EmailTemplates) GetPasswordResetEmailTemplate(baseURL string, verificationToken string) string {
	sanitizedBaseURL := e.normalizeBaseURL(baseURL)
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", sanitizedBaseURL, verificationToken)
	if sanitizedBaseURL == "" {
		resetURL = fmt.Sprintf("/reset-password?token=%s", verificationToken)
	}

	mainContent := fmt.Sprintf(`
		<h1 style="margin:0 0 16px 0; font-size:28px; font-weight:700; color:#0f172a;">Reset your password</h1>
		<p style="margin:0 0 24px 0; font-size:16px; line-height:1.6; color:#1f2937;">
			Hello! You can reset your password by using the secure button below.
		</p>
		<table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 auto 24px auto;">
			<tr>
				<td align="center" style="border-radius:999px; background-color:#2563eb;">
					<a href="%s" style="display:inline-block; padding:14px 32px; font-size:16px; font-weight:600; color:#ffffff; text-decoration:none; border-radius:999px;">Reset Password</a>
				</td>
			</tr>
		</table>
		<p style="margin:0 0 8px 0; font-size:14px; line-height:1.6; color:#475569;">
			Or copy and paste this link into your browser:
		</p>
		<p style="margin:0 0 24px 0; font-size:14px; line-height:1.6;">
			<a href="%s" style="color:#2563eb; word-break:break-all;">%s</a>
		</p>
		<p style="margin:0; font-size:14px; line-height:1.6; color:#475569;">
			If you did not request this email, please ignore it.
		</p>
	`, resetURL, resetURL, resetURL)

	return e.wrapEmail(sanitizedBaseURL, mainContent)
}

func (e *EmailTemplates) wrapEmail(baseURL string, mainContent string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Smarter Lynx</title>
</head>
<body style="margin:0; padding:0; background-color:#f5f7fb; font-family:'Helvetica Neue', Arial, sans-serif; color:#0f172a;">
	<table role="presentation" width="100%%" cellpadding="0" cellspacing="0">
		<tr>
			<td align="center" style="padding:32px 16px;">
				<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px; background-color:#ffffff; border-radius:16px; overflow:hidden; border:1px solid #e2e8f0; box-shadow:0 12px 30px rgba(15, 23, 42, 0.08);">
					<tr>
						<td style="padding:40px 32px;">
							%s
						</td>
					</tr>
				</table>
				%s
			</td>
		</tr>
	</table>
</body>
</html>
`, mainContent, e.buildFooter(baseURL))
}

func (e *EmailTemplates) buildFooter(baseURL string) string {
	sanitized := e.normalizeBaseURL(baseURL)
	websiteURL := sanitized
	if websiteURL == "" {
		websiteURL = defaultWebsiteURL
	}

	contactBase := strings.TrimRight(websiteURL, "/")
	if contactBase == "" {
		contactBase = defaultWebsiteURL
	}
	contactURL := fmt.Sprintf("%s/contact", contactBase)

	return fmt.Sprintf(`
	<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px; margin-top:24px;">
		<tr>
			<td style="text-align:center; color:#64748b; font-size:13px; line-height:1.7;">
				<p style="margin:0 0 8px 0; font-weight:600; color:#475569;">Smarter Lynx</p>
				<p style="margin:0 0 12px 0;">
					<a href="%s" style="color:#2563eb; text-decoration:none; font-weight:600;">Website</a>
					<span style="margin:0 10px; color:#cbd5f5;">•</span>
					<a href="%s" style="color:#2563eb; text-decoration:none; font-weight:600;">Contact Support</a>
					<span style="margin:0 10px; color:#cbd5f5;">•</span>
					<a href="%s" style="color:#2563eb; text-decoration:none; font-weight:600;">Facebook</a>
				</p>
				<p style="margin:0;">You're receiving this email because you requested it from Smarter Lynx.</p>
			</td>
		</tr>
	</table>
`, websiteURL, contactURL, defaultFacebookURL)
}

func (e *EmailTemplates) normalizeBaseURL(baseURL string) string {
	trimmed := strings.TrimSpace(baseURL)
	return strings.TrimRight(trimmed, "/")
}
