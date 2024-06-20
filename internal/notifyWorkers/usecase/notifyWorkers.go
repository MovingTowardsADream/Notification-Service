package usecase

import (
	"Notification_Service/internal/entity"
	"context"
	"fmt"
	"time"
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
	time.Sleep(1 * time.Second)
	fmt.Println("SEND MESSAGE ON MAIL: ", notify)

	return nil
}

func (uc *NotifyWorkerUseCase) CreateNewPhoneNotify(ctx context.Context, notify entity.PhoneDate) error {
	time.Sleep(1 * time.Second)
	fmt.Println("SEND MESSAGE ON PHONE: ", notify)

	return nil
}
