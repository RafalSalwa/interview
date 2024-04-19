package email

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Config struct {
	Host string
	Port int
	From string
}

type Client struct {
	client      *mail.SMTPServer
	DefaultFrom string
}

func NewClient(cfg Config) Client {
	client := mail.NewSMTPClient()
	client.Host = cfg.Host
	client.Port = cfg.Port
	email := Client{
		client:      client,
		DefaultFrom: cfg.From,
	}
	return email
}

type UserEmailData struct {
	Username         string
	Email            string
	VerificationCode string
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

	tmpl, err := ParseTemplateDir("templates")
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
	defer func(con *mail.SMTPClient) error {
		err = con.Close()
		if err != nil {
			return err
		}
		return nil
	}(con)
	if err != nil {
		return err
	}

	err = m.Send(con)
	if err != nil {
		return err
	}

	return nil
}
