package store

import (
	"context"
	"database/sql"
	"email-sender/helpers"
	"email-sender/models"
	"email-sender/pkg/logger"
	"github.com/lib/pq"
)

type MessageStore struct {
	db helpers.DbConnection
	logger.Logger
	smpt *helpers.SmtpHandler
}

func NewMessageStore(db helpers.DbConnection, lg logger.Logger, smtp *helpers.SmtpHandler) *MessageStore {
	return &MessageStore{
		db:     db,
		Logger: lg,
		smpt:   smtp,
	}
}

func (ms MessageStore) MessageHanding(message models.MessageRequest) error {
	ms.Infof("Message from %s received", message.Sender)

	err := ms.InsertMessageRequest(context.Background(), message)
	if err != nil {
		ms.Errorf("Message from %s save in db error %v", err, message.Sender)
		return err
	}

	ms.Infof("Message from %s saved in db", message.Sender)

	err = ms.smpt.SendMail(&message)
	if err != nil {
		ms.Errorf("Message send error %v", err)
		return err
	}

	ms.Infof("Message from %s send success", message.Sender)

	message.State = StateSendSuccess

	err = ms.UpdateStateMessageRequest(context.Background(), message)
	if err != nil {
		ms.Errorf("Message status in db change error %v", err)
	}

	return nil
}

func (ms MessageStore) InsertMessageRequest(ctx context.Context, message models.MessageRequest) error {
	messages := message.ConvertToMessage()
	for _, v := range messages {
		err := ms.insert(ctx, v)
		if err != nil {
			ms.errorHandle(v, err)
			return err
		}

	}
	return nil
}

func (ms MessageStore) errorHandle(message models.Message, err error) {
	pqErr, ok := err.(*pq.Error)
	if ok {
		if pqErr.Code == "42P01" {
			ms.db.AutoMigrate(&message)
		}
	}
}

func (ms MessageStore) UpdateStateMessageRequest(ctx context.Context, message models.MessageRequest) error {
	messages := message.ConvertToMessage()
	for _, v := range messages {
		err := ms.updateState(ctx, v)
		if err != nil {
			ms.errorHandle(v, err)
			return err
		}

	}
	return nil
}

func (ms MessageStore) updateState(ctx context.Context, message models.Message) error {
	tx := ms.db.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Model(&message).Updates(models.Message{State: message.State}).Where("unique_id = ?", message.UniqueId).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (ms MessageStore) insert(ctx context.Context, message models.Message) error {
	tx := ms.db.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}
