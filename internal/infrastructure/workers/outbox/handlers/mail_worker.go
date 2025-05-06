//nolint:dupl // will be rewritten to interface later
package handlers

import (
	"context"
	"errors"
	"time"

	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/logger"
)

const (
	defaultMailBatchSize = 10
	defaultMailTimeout   = 1 * time.Second
)

var ErrMailWorkerStop = errors.New("mail worker stopping")

type MailData interface {
	GetBatchMailNotify(ctx context.Context, batch *dto.BatchNotify) ([]*dto.MailIdempotencyData, error)
	ProcessedBatchMailNotify(ctx context.Context, keys []*dto.IdempotencyKey) error
}

type MailGateway interface {
	CreateMailNotify(ctx context.Context, notify *dto.MailIdempotencyData) error
}

type MailOption func(worker *MailWorker)

func MailBatchSize(batchSize uint64) MailOption {
	return func(w *MailWorker) {
		w.batchSize = batchSize
	}
}

func MailTimeout(timeout time.Duration) MailOption {
	return func(w *MailWorker) {
		w.timeout = timeout
	}
}

type MailWorker struct {
	log      logger.Logger
	mailData MailData
	gateway  MailGateway

	batchSize uint64
	timeout   time.Duration
	stop      chan struct{}
}

func NewMailWorker(log logger.Logger, mailData MailData, gateway MailGateway, opts ...MailOption) *MailWorker {
	mw := &MailWorker{
		log:       log,
		mailData:  mailData,
		gateway:   gateway,
		batchSize: defaultMailBatchSize,
		timeout:   defaultMailTimeout,
		stop:      make(chan struct{}),
	}

	for _, opt := range opts {
		opt(mw)
	}

	return mw
}

func (w *MailWorker) Run() error {
	const op = "MailWorker.Run"

	for {
		select {
		case <-w.stop:
			return ErrMailWorkerStop
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

	mailData, err := w.mailData.GetBatchMailNotify(context.Background(), &dto.BatchNotify{BatchSize: w.batchSize})

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

	err = w.mailData.ProcessedBatchMailNotify(context.Background(), keys)

	if err != nil {
		w.log.Error(op, w.log.Err(err))
		return err
	}

	return nil
}

func (w *MailWorker) Stop() {
	close(w.stop)
}
