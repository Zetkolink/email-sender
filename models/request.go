package models

import "time"

type MessageRequest struct {
	Sender  string   `json:"sender"`
	To      []string `json:"to"`
	Subject *string  `json:"subject"`
	Message string   `json:"message"`
}

func (m MessageRequest) ConvertToMessage() []Message {
	messages := make([]Message, 0, len(m.To))
	for _, v := range m.To {
		message := Message{
			Sender:    m.Sender,
			To:        v,
			Subject:   m.Subject,
			Message:   m.Message,
			CreatedAt: time.Now(),
		}
		messages = append(messages, message)
	}
	return messages
}
