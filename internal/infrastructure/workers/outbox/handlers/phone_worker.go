package handlers

import (
	"context"
	"errors"
	"time"

	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/logger"
)

type PhoneData interface {
	GetBatchPhoneNotify(ctx context.Context, batch *dto.BatchNotify) ([]*dto.PhoneIdempotencyData, error)
	ProcessedBatchPhoneNotify(ctx context.Context, keys []*dto.IdempotencyKey) error
}

type PhoneGateway interface {
	CreatePhoneNotify(ctx context.Context, notify *dto.PhoneIdempotencyData) error
}

type PhoneWorker struct {
	log       logger.Logger
	phoneData PhoneData
	gateway   PhoneGateway

	ctx       context.Context
	batchSize uint64
	timeout   time.Duration
	stop      chan struct{}
}

func NewPhoneWorker(log logger.Logger, phoneData PhoneData, gateway PhoneGateway, batchSize uint64, timeout time.Duration) *PhoneWorker {
	return &PhoneWorker{
		log:       log,
		phoneData: phoneData,
		gateway:   gateway,
		ctx:       context.Background(),
		batchSize: batchSize,
		timeout:   timeout,
		stop:      make(chan struct{}),
	}
}

func (w *PhoneWorker) Run() error {
	const op = "PhoneWorker.Run"

	for {
		select {
		case <-w.stop:
			return errors.New("phone worker stopping")
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

	phoneData, err := w.phoneData.GetBatchPhoneNotify(w.ctx, &dto.BatchNotify{BatchSize: w.batchSize})

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

	err = w.phoneData.ProcessedBatchPhoneNotify(w.ctx, keys)

	if err != nil {
		w.log.Error(op, w.log.Err(err))
		return err
	}

	return nil
}
