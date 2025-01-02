package smtp

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

const (
	defaultSMTPPort = 443
)

type Params struct {
	Domain, Username, Password, Mail string
}

type SMTP struct {
	Dialer *gomail.Dialer

	Mail string
	Port int
}

func New(params *Params, opts ...Option) *SMTP {
	smtp := &SMTP{
		Port: defaultSMTPPort,
		Mail: params.Mail,
	}

	for _, opts := range opts {
		opts(smtp)
	}

	dialer := gomail.NewDialer(params.Domain, smtp.Port, params.Username, params.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtp.Dialer = dialer

	return smtp
}
