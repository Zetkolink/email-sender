package models

import "time"

// MessageRequest структура сообщений.
type MessageRequest struct {
	// Sender отправитель.
	Sender string `json:"sender"`

	// UniqueId уникальный ID сообщения.
	UniqueId string `json:"unique_id"`

	// To массив получателей.
	To []string `json:"to"`

	// Subject тема сообщения.
	Subject *string `json:"subject"`

	// Message текст сообщения.
	Message string `json:"message"`

	// State статус сообщения.
	State string `json:"state"`
}

// ConvertToMessage конвертация сообщений для хранения их в базе в соответствии первой нормальной формы.
func (m MessageRequest) ConvertToMessage() []Message {
	messages := make([]Message, 0, len(m.To))
	tt := time.Now().Format(time.RFC3339)
	for _, v := range m.To {
		message := Message{
			Sender:    m.Sender,
			To:        v,
			UniqueId:  m.UniqueId,
			Subject:   m.Subject,
			Message:   m.Message,
			CreatedAt: tt,
			State:     false,
		}
		messages = append(messages, message)
	}
	return messages
}
