package smtp

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

const (
	defaultSMTPPort = 443
)

type Params struct {
	Domain, Username, Password string
}

type SMTP struct {
	Dialer *gomail.Dialer
	Params Params
	Port   int
}

func New(params Params, opts ...Option) *SMTP {
	smtp := &SMTP{
		Port:   defaultSMTPPort,
		Params: params,
	}

	for _, opts := range opts {
		opts(smtp)
	}

	dialer := gomail.NewDialer(params.Domain, smtp.Port, params.Username, params.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // insecure

	smtp.Dialer = dialer

	return smtp
}
