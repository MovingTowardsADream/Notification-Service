package smtp

import (
	"context"
	"fmt"

	//"gopkg.in/gomail.v2"

	"Notification_Service/internal/interfaces/dto"
)

type WorkerMail struct {
	sender *SMTP
}

func NewWorkerMail(smtp *SMTP) *WorkerMail {
	return &WorkerMail{sender: smtp}
}

func (uc *WorkerMail) SendMailLetter(_ context.Context, notify dto.MailDate) error {
	//const op = "smtp - SendMailLetter"
	//
	//m := gomail.NewMessage()
	//
	//m.SetHeader("From", uc.sender.Params.Username)
	//m.SetHeader("To", notify.Mail)
	//m.SetHeader("Subject", notify.Subject)
	//m.SetBody("text/html", notify.Body)
	//
	//if err := uc.sender.Dialer.DialAndSend(m); err != nil {
	//	return fmt.Errorf("%s - uc.sender.Dialer.DialAndSend: %w", op, err)
	//}

	fmt.Println("SEND MESSAGE ON MAIL: ", notify)

	return nil
}
