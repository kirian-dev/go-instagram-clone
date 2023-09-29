package helpers

import (
	"bytes"
	"crypto/tls"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	"go-instagram-clone/services/chat/internal/domain/models"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func ParseTemplateDir(relativeDir string) (*template.Template, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(currentDir, relativeDir)

	var paths []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

func SendEmail(user *models.UserResponse, data *EmailData, emailTemp string, cfg *config.Config, log *logger.ZapLogger, imageByte []byte, imageBase64 string, imageName string) {
	from := cfg.EmailFrom
	smtpPass := cfg.SMTPPassword
	smtpUser := cfg.SMTPUser
	to := user.Email
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, emailTemp, &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)

	if imageByte != nil && imageBase64 != "" {
		m.Embed(imageName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(imageByte)
			return err
		}), gomail.SetHeader(map[string][]string{
			"Content-Type":              {"image/png"},
			"Content-ID":                {"<" + imageName + ">"},
			"Content-Transfer-Encoding": {"base64"},
		}))
	}

	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, utils.ParsePort(smtpPort), smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}
