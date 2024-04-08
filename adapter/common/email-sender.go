package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Sender struct {
	auth       smtp.Auth
	host       string
	userName   string
	password   string
	portNumber string
}

type Message struct {
	To           []string
	Cc           []string
	Bcc          []string
	subject      string
	body         string
	attachments  map[string][]byte
	htmlTemplate *template.Template
	data         any
}

func NewSender(host string, userName string, password string, portNumber string) *Sender {
	auth := smtp.PlainAuth("", userName, password, host)
	return &Sender{auth, host, userName, password, portNumber}
}

func (s *Sender) Send(message *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", s.host, s.portNumber), s.auth, s.userName, message.To, message.ToBytes())
}

func NewTextMessage(subject, body string) *Message {
	return &Message{subject: subject, body: body, attachments: make(map[string][]byte), htmlTemplate: nil}
}

func NewHtmlMessage(subject string, data any, templateFileName string) (*Message, error) {
	var err error

	htmlTemplate, err := template.ParseFiles(templateFileName)

	if err != nil {
		return &Message{}, nil
	}

	return &Message{subject: subject, attachments: make(map[string][]byte), htmlTemplate: htmlTemplate, data: data}, nil
}

func (m *Message) AddAttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.Cc, ",")))
	}

	if len(m.Bcc) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.Bcc, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}

	if m.htmlTemplate != nil {
		buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\n")
		m.htmlTemplate.Execute(buf, m.data)
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
		buf.WriteString(m.body)
	}

	if withAttachments {
		for k, v := range m.attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
