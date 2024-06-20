package usecase

import (
	"Notification_Service/internal/entity"
	"context"
	"fmt"
)

type NotifyWorker interface {
	CreateNewMailNotify(ctx context.Context, notify entity.MailDate) error
	CreateNewPhoneNotify(ctx context.Context, notify entity.PhoneDate) error
}

type NotifyWorkerUseCase struct {
}

func NewNotifyWorker() *NotifyWorkerUseCase {
	return &NotifyWorkerUseCase{}
}

func (uc *NotifyWorkerUseCase) CreateNewMailNotify(ctx context.Context, notify entity.MailDate) error {

	// TODO Send message on mail

	fmt.Println("SEND MESSAGE ON MAIL: ", notify)

	return nil
}

func (uc *NotifyWorkerUseCase) CreateNewPhoneNotify(ctx context.Context, notify entity.PhoneDate) error {

	// TODO Send message on phone

	fmt.Println("SEND MESSAGE ON PHONE: ", notify)

	return nil
}
