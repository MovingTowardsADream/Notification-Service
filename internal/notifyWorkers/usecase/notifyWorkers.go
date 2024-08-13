package usecase

import (
	"Notification_Service/internal/entity"
	"context"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
)

type NotifyWorker interface {
	CreateNewMailNotify(ctx context.Context, notify entity.MailDate) error
	CreateNewPhoneNotify(ctx context.Context, notify entity.PhoneDate) error
}

type SMTP struct {
	Domain   string
	Port     int
	UserName string
	Password string
}

type NotifyWorkerUseCase struct {
	mail   string
	dialer *gomail.Dialer
}

func NewNotifyWorker(smtp SMTP, mail string) *NotifyWorkerUseCase {
	d := gomail.NewDialer(smtp.Domain, smtp.Port, smtp.UserName, smtp.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &NotifyWorkerUseCase{mail: mail, dialer: d}
}

func (uc *NotifyWorkerUseCase) CreateNewMailNotify(ctx context.Context, notify entity.MailDate) error {
	//m := gomail.NewMessage()
	//m.SetHeader("From", uc.mail)
	//m.SetHeader("To", notify.Mail)
	//m.SetHeader("Subject", notify.Subject)
	//m.SetBody("text/html", notify.Body)
	//
	//if err := uc.dialer.DialAndSend(m); err != nil {
	//	return err
	//}

	// TODO Send message on email

	fmt.Println("SEND MESSAGE ON PHONE: ", notify)

	return nil
}

func (uc *NotifyWorkerUseCase) CreateNewPhoneNotify(ctx context.Context, notify entity.PhoneDate) error {

	// TODO Send message on phone

	fmt.Println("SEND MESSAGE ON PHONE: ", notify)

	return nil
}
