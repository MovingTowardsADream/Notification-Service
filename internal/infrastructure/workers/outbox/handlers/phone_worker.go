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
	defaultPhoneBatchSize = 10
	defaultPhoneTimeout   = 1 * time.Second
)

var ErrPhoneWorkerStop = errors.New("phone worker stopping")

type PhoneData interface {
	GetBatchPhoneNotify(ctx context.Context, batch *dto.BatchNotify) ([]*dto.PhoneIdempotencyData, error)
	ProcessedBatchPhoneNotify(ctx context.Context, keys []*dto.IdempotencyKey) error
}

type PhoneGateway interface {
	CreatePhoneNotify(ctx context.Context, notify *dto.PhoneIdempotencyData) error
}

type PhoneOption func(worker *PhoneWorker)

func PhoneBatchSize(batchSize uint64) PhoneOption {
	return func(w *PhoneWorker) {
		w.batchSize = batchSize
	}
}

func PhoneTimeout(timeout time.Duration) PhoneOption {
	return func(w *PhoneWorker) {
		w.timeout = timeout
	}
}

type PhoneWorker struct {
	log       logger.Logger
	phoneData PhoneData
	gateway   PhoneGateway

	batchSize uint64
	timeout   time.Duration
	stop      chan struct{}
}

func NewPhoneWorker(log logger.Logger, phoneData PhoneData, gateway PhoneGateway, opts ...PhoneOption) *PhoneWorker {
	pw := &PhoneWorker{
		log:       log,
		phoneData: phoneData,
		gateway:   gateway,
		batchSize: defaultPhoneBatchSize,
		timeout:   defaultPhoneTimeout,
		stop:      make(chan struct{}),
	}

	for _, opt := range opts {
		opt(pw)
	}

	return pw
}

func (w *PhoneWorker) Run() error {
	const op = "PhoneWorker.Run"

	for {
		select {
		case <-w.stop:
			return ErrPhoneWorkerStop
		default:
			if err := w.worker(); err != nil {
				w.log.Error(op, w.log.Err(err))
			}
		}
		time.Sleep(w.timeout)
	}
}

func (w *PhoneWorker) worker() error {
	const op = "PhoneWorker.worker()"

	phoneData, err := w.phoneData.GetBatchPhoneNotify(context.Background(), &dto.BatchNotify{BatchSize: w.batchSize})

	if err != nil {
		w.log.Error(op, w.log.Err(err))
		return err
	}

	keys := make([]*dto.IdempotencyKey, 0, len(phoneData))

	for _, phone := range phoneData {
		err := w.gateway.CreatePhoneNotify(context.Background(), phone)
		if err != nil {
			w.log.Error(op, w.log.Err(err))
			return err
		}
		keys = append(keys, &dto.IdempotencyKey{RequestID: phone.RequestID})
	}

	err = w.phoneData.ProcessedBatchPhoneNotify(context.Background(), keys)

	if err != nil {
		w.log.Error(op, w.log.Err(err))
		return err
	}

	return nil
}

func (w *PhoneWorker) Stop() {
	close(w.stop)
}
