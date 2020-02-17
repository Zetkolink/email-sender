package models

// Message тело запроса на отправку email.
type Message struct {
	// ID идентификатор сообщения.
	ID uint `gorm:"primary_key"`

	// UniqueId уникальный ID сообщения.
	UniqueId string `json:"unique_id"`

	// CreatedAt дата создания.
	CreatedAt string

	// Sender отправитель.
	Sender string `json:"sender"`

	// To получатель.
	To string `json:"to"`

	// Subject тема сообщения.
	Subject *string `json:"subject"`

	// Message текст сообщения.
	Message string `json:"message"`

	// State статус доставки.
	State bool `json:"state"`
}
