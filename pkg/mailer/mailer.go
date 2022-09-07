package mailer

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Mail struct {
	from    string
	To      []string
	Subject string
	Body    string
}

func (m *Mail) Build() []byte {
	var buff = bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("Content-Type: %s\r\n", gin.MIMEHTML))
	buff.WriteString(fmt.Sprintf("From: %s\r\n", m.from))
	if len(m.To) > 0 {
		buff.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ";")))
	}

	buff.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	buff.WriteString("\r\n\n" + m.Body)

	return buff.Bytes()
}

type Mailer interface {
	Send(mail Mail) error
}
