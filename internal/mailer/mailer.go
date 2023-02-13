package mailer

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (Mailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		return Mailer{}, err
	}

	mailer := Mailer{
		dialer: client,
		sender: sender,
	}

	return mailer, nil
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()

	if err := msg.To(recipient); err != nil {
		return err
	}

	if err := msg.From(m.sender); err != nil {
		return err
	}

	msg.Subject(subject.String())

	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.SetBodyString(mail.TypeTextHTML, htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
