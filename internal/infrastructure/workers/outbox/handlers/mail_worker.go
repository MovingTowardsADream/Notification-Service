package handlers

import (
	"context"
	"errors"
	"time"

	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/logger"
)

type MailData interface {
	GetBatchMailNotify(ctx context.Context, batch *dto.BatchNotify) ([]*dto.MailIdempotencyData, error)
	ProcessedBatchMailNotify(ctx context.Context, keys []*dto.IdempotencyKey) error
}

type MailGateway interface {
	CreateMailNotify(ctx context.Context, notify *dto.MailIdempotencyData) error
}

type MailWorker struct {
	log      logger.Logger
	mailData MailData
	gateway  MailGateway

	ctx       context.Context
	batchSize uint64
	timeout   time.Duration
	stop      chan struct{}
}

func NewMailWorker(log logger.Logger, mailData MailData, gateway MailGateway, batchSize uint64, timeout time.Duration) *MailWorker {
	return &MailWorker{
		log:       log,
		mailData:  mailData,
		gateway:   gateway,
		ctx:       context.Background(),
		batchSize: batchSize,
		timeout:   timeout,
		stop:      make(chan struct{}),
	}
}

func (w *MailWorker) Run() error {
	const op = "MailWorker.Run"

	for {
		select {
		case <-w.stop:
			return errors.New("mail worker stopping")
		default:
			if err := w.worker(); err != nil {
				w.log.Error(op, w.log.Err(err))
			}
		}
		time.Sleep(w.timeout)
	}
}

func (w *MailWorker) worker() error {
	const op = "MailWorker.worker()"

	mailData, err := w.mailData.GetBatchMailNotify(w.ctx, &dto.BatchNotify{BatchSize: w.batchSize})

	if err != nil {
		w.log.Error(op, w.log.Err(err))
		return err
	}

	keys := make([]*dto.IdempotencyKey, 0, len(mailData))

	for _, mail := range mailData {
		err := w.gateway.CreateMailNotify(context.Background(), mail)
		if err != nil {
			w.log.Error(op, w.log.Err(err))
			return err
		}
		keys = append(keys, &dto.IdempotencyKey{RequestID: mail.RequestID})
	}

	err = w.mailData.ProcessedBatchMailNotify(w.ctx, keys)

	if err != nil {
		w.log.Error(op, w.log.Err(err))
		return err
	}

	return nil
}
