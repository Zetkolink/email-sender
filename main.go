package main

import (
	"email-sender/helpers"
	"email-sender/models"
	"email-sender/pkg/logger"
	"email-sender/store"
	"encoding/json"
	"os"
	"sync"
)

func main() {
	cfg := helpers.InitConfig()

	db := helpers.InitDb(cfg)

	rb := helpers.InitRabbit(cfg)

	smtp := helpers.InitSmtp(cfg)

	lg := logger.New(os.Stderr, cfg.LogLevel, cfg.LogFormat)

	ms := store.NewMessageStore(db, lg, smtp)

	consumeCh := rb.InitChannel(cfg)

	go func() {
		d, err := consumeCh.Consume(cfg.Rb.Queue, "", false, false, false, false, nil)
		if err != nil {
			lg.Errorf("Consume msg error %v", err)
		}

		for msg := range d {
			message := models.MessageRequest{}
			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				ms.Errorf("Message unmarshal error %v", err)
				_ = msg.Ack(false)
			}
			err = ms.MessageHanding(message)
			if err != nil {
				_ = msg.Nack(false, true)
			}
			_ = msg.Ack(false)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	wg.Wait()
}
