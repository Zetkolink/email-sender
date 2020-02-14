package store

import (
	"context"
	"database/sql"
	"email-sender/helpers"
	"email-sender/models"
	"email-sender/pkg/logger"
)

type MessageStore struct {
	db helpers.DbConnection
	logger.Logger
}

func NewMessageStore(db helpers.DbConnection, lg logger.Logger) *MessageStore {
	return &MessageStore{
		db:     db,
		Logger: lg,
	}
}

func (ms MessageStore) Insert(ctx context.Context, message models.Message) {
	tx := ms.db.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		ms.Errorf("Insert message begin transaction error %s", err)
	}

	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		ms.Errorf("Insert message create error %s", err)
	}

	err := tx.Commit().Error
	if err != nil {
		ms.Errorf("Insert message commit transaction error %s", err)
	}
}
