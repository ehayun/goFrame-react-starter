package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type EmailService struct {
	host         string
	port         int
	username     string
	password     string
	from         string
	fromName     string
	templatePath string
}

func NewEmailService() *EmailService {
	ctx := gctx.New()
	cfg := g.Cfg()

	return &EmailService{
		host:         cfg.MustGet(ctx, "email.smtp.host").String(),
		port:         cfg.MustGet(ctx, "email.smtp.port").Int(),
		username:     os.Getenv("EMAIL_USERNAME"),
		password:     os.Getenv("EMAIL_PASSWORD"),
		from:         os.Getenv("EMAIL_FROM"),
		fromName:     cfg.MustGet(ctx, "email.smtp.fromName").String(),
		templatePath: "templates/email",
	}
}

type EmailData struct {
	To      []string
	Subject string
	Body    string
}

func (es *EmailService) Send(data *EmailData) error {
	// Authentication
	auth := smtp.PlainAuth("", es.username, es.password, es.host)

	// Build message
	msg := es.buildMessage(data)

	// Send email
	addr := fmt.Sprintf("%s:%d", es.host, es.port)
	if err := smtp.SendMail(addr, auth, es.from, data.To, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (es *EmailService) buildMessage(data *EmailData) []byte {
	msg := fmt.Sprintf("From: %s <%s>\r\n", es.fromName, es.from)
	msg += fmt.Sprintf("To: %s\r\n", data.To[0])
	msg += fmt.Sprintf("Subject: %s\r\n", data.Subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += data.Body

	return []byte(msg)
}

type WelcomeEmailData struct {
	Name string
	Link string
}

func (es *EmailService) SendWelcomeEmail(to, name string) error {
	// Parse template
	tmplPath := filepath.Join(es.templatePath, "welcome.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var body bytes.Buffer
	data := WelcomeEmailData{
		Name: name,
		Link: "https://tzlev.com/dashboard",
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Send email
	return es.Send(&EmailData{
		To:      []string{to},
		Subject: "Welcome to Tzlev!",
		Body:    body.String(),
	})
}
