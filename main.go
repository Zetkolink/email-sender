package main

import (
	"email-sender/helpers"
	"email-sender/pkg/logger"
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	cfg := helpers.InitConfig()

	//db := helpers.InitDb(cfg)

	rb := helpers.InitRabbit(cfg)

	//smtp := helpers.InitSmtp(cfg)

	lg := logger.New(os.Stderr, cfg.LogLevel, cfg.LogFormat)

	sendCh, err := rb.Channel()
	if err != nil {
		log.Panic(err)
	}

	exchangeName := "test-exchange"

	err = sendCh.ExchangeDeclare(exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Panic(err)
	}

	_, err = sendCh.QueueDeclare(cfg.Rb.Channel, true, false, false, false, nil)
	if err != nil {
		log.Panic(err)
	}

	if err := sendCh.QueueBind(cfg.Rb.Channel, "", exchangeName, false, nil); err != nil {
		log.Panic(err)
	}

	go func() {
		for {
			err := sendCh.Publish(exchangeName, "", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(time.Now().String()),
			})
			log.Printf("publish, err: %v", err)
			time.Sleep(5 * time.Second)
			err = sendCh.ExchangeDeclare(exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
			if err != nil {
				//log.Panic(err)
			}

			_, err = sendCh.QueueDeclare(cfg.Rb.Channel, true, false, false, false, nil)
			if err != nil {
				//log.Panic(err)
			}

			if err := sendCh.QueueBind(cfg.Rb.Channel, "", exchangeName, false, nil); err != nil {
				//log.Panic(err)
			}
		}
	}()

	consumeCh, err := rb.Channel()
	if err != nil {
		lg.Fatalf("init channel failed")
	}

	go func() {
		d, err := consumeCh.Consume(cfg.Rb.Channel, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err)
		}

		for msg := range d {
			log.Printf("msg1: %s", string(msg.Body))
			_ = msg.Ack(true)
		}
	}()

	go func() {
		d, err := consumeCh.Consume(cfg.Rb.Channel, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err)
		}

		for msg := range d {
			log.Printf("msg2: %s", string(msg.Body))
			_ = msg.Nack(false, true)
		}
	}()

	go func() {
		d, err := consumeCh.Consume(cfg.Rb.Channel, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err)
		}

		for msg := range d {
			log.Printf("msg3: %s", string(msg.Body))
			_ = msg.Nack(false, false)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	wg.Wait()

	//ms := store.NewMessageStore(db, lg)

}
