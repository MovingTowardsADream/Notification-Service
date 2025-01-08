package smtp

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"

	"Notification_Service/internal/interfaces/dto"
)

type NotifyWorker interface {
	CreateNewMailNotify(ctx context.Context, notify dto.MailDate) error
	CreateNewPhoneNotify(ctx context.Context, notify dto.PhoneDate) error
}

type NotifyWorkerUseCase struct {
	SMTP *SMTP
}

func NewNotifyWorker(smtp *SMTP) *NotifyWorkerUseCase {
	return &NotifyWorkerUseCase{SMTP: smtp}
}

func (uc *NotifyWorkerUseCase) CreateNewMailNotify(ctx context.Context, notify dto.MailDate) error {
	m := gomail.NewMessage()

	m.SetHeader("From", uc.SMTP.Params.Username)
	m.SetHeader("To", notify.Mail)
	m.SetHeader("Subject", notify.Subject)
	m.SetBody("text/html", notify.Body)

	if err := uc.SMTP.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("smtp - NotifyWorkerUseCase - CreateNewMailNotify - uc.SMTP.Dialer.DialAndSend: %w", err)
	}

	return nil
}

func (uc *NotifyWorkerUseCase) CreateNewPhoneNotify(ctx context.Context, notify dto.PhoneDate) error {
	// TODO Send message on phone

	fmt.Println("SEND MESSAGE ON PHONE: ", notify)

	return nil
}
