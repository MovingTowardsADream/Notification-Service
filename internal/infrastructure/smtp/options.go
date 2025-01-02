package smtp

type Option func(*SMTP)

func Port(port int) Option {
	return func(smtp *SMTP) {
		smtp.Port = port
	}
}
