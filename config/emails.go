package config

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"
	mail "gopkg.in/gomail.v2"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

type EmailTypeChoice string

const (
	ET_ACTIVATE              EmailTypeChoice = "activate"
	ET_WELCOME               EmailTypeChoice = "welcome"
	ET_RESET                 EmailTypeChoice = "reset"
	ET_RESET_SUCC            EmailTypeChoice = "reset-success"
	ET_PAYMENT_SUCC          EmailTypeChoice = "payment-succeeded"
	ET_PAYMENT_FAIL          EmailTypeChoice = "payment-failed"
	ET_PAYMENT_CANCEL        EmailTypeChoice = "payment-canceled"
)

func sortEmail(emailType EmailTypeChoice, otp *uint32) map[string]interface{} {
	templateFile := "templates/welcome.html"
	subject := "Account verified"
	data := make(map[string]interface{})
	data["template_file"] = templateFile
	data["subject"] = subject
	data["text"] = "Your Verification was completed."

	// Sort different templates and subject for respective email types
	switch emailType {
	case ET_ACTIVATE:
		templateFile = "templates/email-activation.html"
		subject = "Activate your account"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["otp"] = otp

	case ET_RESET:
		templateFile = "templates/password-reset.html"
		subject = "Activate your account"
		data["template_file"] = templateFile
		data["subject"] = subject
		data["otp"] = otp

	case ET_RESET_SUCC:
		templateFile = "templates/password-reset-success.html"
		subject = "Password reset successfully"
		data["template_file"] = templateFile
		data["subject"] = subject
	}
	return data
}

type EmailContext struct {
	Name string
	Otp *uint32
}

func SendEmail(user *ent.User, emailType EmailTypeChoice, otp *uint32) {
	if os.Getenv("ENVIRONMENT") == "test" {
		return
	}
	cfg := GetConfig()
	emailData := sortEmail(emailType, otp)
	templateFile := emailData["template_file"]
	subject := emailData["subject"]

	// Create a context with dynamic data
	data := EmailContext{
		Name: user.Name,
	}
	if otp, ok := emailData["otp"]; ok {
		otp := otp.(*uint32)
		data.Otp = otp
	}

	// Read the HTML file content
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Unable to identify current directory (needed to load templates)", os.Stderr)
		os.Exit(1)
	}
	basepath := filepath.Dir(file)
	tempfile := fmt.Sprintf("../%s", templateFile.(string))
	htmlContent, err := os.ReadFile(filepath.Join(basepath, tempfile))
	if err != nil {
		log.Fatal("Error reading HTML file:", err)
	}

	// Create a new template from the HTML file content
	tmpl, err := template.New("email_template").Parse(string(htmlContent))
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// Execute the template with the context and set it as the body of the email
	var bodyContent bytes.Buffer
	if err := tmpl.Execute(&bodyContent, data); err != nil {
		log.Fatal("Error executing template:", err)
	}

	// Create a new message
	m := mail.NewMessage()
	m.SetHeader("From", cfg.MailFrom)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject.(string))
	m.SetBody("text/html", bodyContent.String())

	// Create a new SMTP client
	d := mail.NewDialer(cfg.MailSenderHost, cfg.MailSenderPort, cfg.MailSenderEmail, cfg.MailSenderPassword)
	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
	}
}