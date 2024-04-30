package email

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	mail "github.com/xhit/go-simple-mail/v2"
)

type (
	Config struct {
		Host        string
		Port        int
		From        string
		TemplateDir string
	}
	Client struct {
		client      *mail.SMTPServer
		DefaultFrom string
		TemplateDir string
	}
	UserEmailData struct {
		Username         string
		Email            string
		VerificationCode string
	}
)

func NewClient(cfg Config) Client {
	client := mail.NewSMTPClient()
	client.Host = cfg.Host
	client.Port = cfg.Port
	email := Client{
		client:      client,
		DefaultFrom: cfg.From,
		TemplateDir: cfg.TemplateDir,
	}
	return email
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func (c *Client) SendVerificationEmail(data UserEmailData) error {
	var body bytes.Buffer

	tmpl, err := ParseTemplateDir(c.TemplateDir)
	if err != nil {
		log.Fatal("Could not parse tmpl", err)
	}

	tmpl = tmpl.Lookup("verification.html")
	err = tmpl.Execute(&body, &data)
	if err != nil {
		return err
	}

	m := mail.NewMSG()
	m.SetFrom(c.DefaultFrom).
		AddTo(data.Email).
		SetSubject("Please verify your account").
		SetBodyData(mail.TextHTML, body.Bytes())

	con, err := c.client.Connect()
	if err != nil {
		return err
	}
	defer func(con *mail.SMTPClient) {
		_ = con.Close()
	}(con)

	err = m.Send(con)
	if err != nil {
		return err
	}

	return nil
}
