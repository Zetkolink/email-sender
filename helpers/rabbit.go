package helpers

import (
	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"log"
)

type Rabbit struct {
	*rabbitmq.Connection
	amqpUrl string
}

func InitRabbit(cfg Config) Rabbit {
	rb := Rabbit{}
	rb.amqpUrl = cfg.Rb.Amqp
	rb.connect()

	return rb
}

func (r *Rabbit) connect() {
	conn, err := rabbitmq.Dial(r.amqpUrl)
	if err != nil {
		log.Fatal(err)
	}

	r.Connection = conn
}
