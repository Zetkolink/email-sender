package models

import (
	"time"
)

// Message тело запроса на отправку email.
type Message struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	Sender    string  `json:"sender"`
	To        string  `json:"to"`
	Subject   *string `json:"subject"`
	Message   string  `json:"message"`
}
