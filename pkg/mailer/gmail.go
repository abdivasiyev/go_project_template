package mailer

import (
	"crypto/tls"
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/smtp"
)

var Module = fx.Provide(NewGmail)

type GmailParams struct {
	fx.In
	Config config.Config
	Log    logger.Logger
	Sentry sentry.Handler
}

type gmailMailer struct {
	config config.Config
	log    logger.Logger
	sentry sentry.Handler
}

func NewGmail(params GmailParams) Mailer {
	return &gmailMailer{
		config: params.Config,
		log:    params.Log,
		sentry: params.Sentry,
	}
}

func (g *gmailMailer) getClient() (*smtp.Client, error) {
	var (
		smtpAddr  = fmt.Sprintf("%s:%d", g.config.GetString(config.SmtpHostKey), g.config.GetInt(config.SmtpPortKey))
		tlsConfig = &tls.Config{ServerName: g.config.GetString(config.SmtpHostKey), InsecureSkipVerify: true}
	)

	conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
	if err != nil {
		g.sentry.HandleError(err)
		g.log.Error("could not dial with smtp server", zap.Error(err))
		return nil, err
	}

	client, err := smtp.NewClient(conn, g.config.GetString(config.SmtpHostKey))
	if err != nil {
		g.sentry.HandleError(err)
		g.log.Error("could not create new client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

func (g *gmailMailer) Send(mail Mail) error {
	var (
		auth = smtp.PlainAuth(
			"",
			g.config.GetString(config.SmtpUsernameKey),
			g.config.GetString(config.SmtpPasswordKey),
			g.config.GetString(config.SmtpHostKey),
		)
	)

	mail.from = g.config.GetString(config.SmtpUsernameKey)

	client, err := g.getClient()
	if err != nil {
		g.sentry.HandleError(err)
		g.log.Error("could not authenticate to smtp server", zap.Error(err))
		return err
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		g.sentry.HandleError(err)
		g.log.Error("could not authenticate to smtp server", zap.Error(err))
		return err
	}

	if err = client.Mail(mail.from); err != nil {
		g.sentry.HandleError(err)
		g.log.Error("could not mail", zap.Error(err))
		return err
	}

	for _, toEmail := range mail.To {
		if err = client.Rcpt(toEmail); err != nil {
			g.sentry.HandleError(err)
			g.log.Error("could not rcpt", zap.Error(err))
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		g.sentry.HandleError(err)
		g.log.Error("could not make DATA command", zap.Error(err))
		return err
	}
	defer writer.Close()

	n, err := writer.Write(mail.Build())
	if err != nil {
		g.sentry.HandleError(err)
		g.log.Error("could not write mail body", zap.Error(err))
		return err
	}

	if n == 0 {
		g.sentry.HandleMessage("nothing was written to mail: %v", mail)
		g.log.Error("nothing was written to mail", zap.Any("mail", mail))
	}

	return nil
}
