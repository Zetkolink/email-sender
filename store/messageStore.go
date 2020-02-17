package store

import (
	"context"
	"database/sql"
	"email-sender/helpers"
	"email-sender/models"
	"email-sender/pkg/logger"
	"github.com/jinzhu/gorm"
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
	ms.Infof("Message UNIQUE_ID:'%s' received", message.UniqueId)

	tx, err := ms.InsertMessageRequest(context.Background(), message)
	if err != nil {
		ms.Errorf("Message UNIQUE_ID:'%s' prepare transaction in db error %v", err, message.UniqueId)
		return err
	}

	ms.Infof("Message UNIQUE_ID:'%s' prepare transaction in db", message.UniqueId)

	err = ms.smpt.SendMail(&message)
	if err != nil {
		tx.Rollback()
		ms.Errorf("Message send error %v", err)
		return err
	}

	ms.Infof("Message UNIQUE_ID:'%s' send success", message.UniqueId)
	err = tx.Commit().Error
	if err != nil {
		ms.Errorf("Message save in db error %v", err)
		return err
	}
	ms.Infof("Message UNIQUE_ID:'%s' transaction commit success", message.UniqueId)

	return nil
}

func (ms MessageStore) InsertMessageRequest(ctx context.Context, message models.MessageRequest) (*gorm.DB, error) {
	messages := message.ConvertToMessage()

	tx := ms.db.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return tx, err
	}

	for _, v := range messages {
		v.State = true
		err := ms.insert(tx, v)
		if err != nil {
			ms.errorHandle(v, err)
			return tx, err
		}

	}
	return tx, nil
}

func (ms MessageStore) errorHandle(message models.Message, err error) {
	pqErr, ok := err.(*pq.Error)
	if ok {
		if pqErr.Code == "42P01" {
			ms.db.AutoMigrate(&message)
		}
	}
}

func (ms MessageStore) insert(tx *gorm.DB, message models.Message) error {
	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
