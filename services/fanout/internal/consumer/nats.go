package consumer

import (
	"encoding/json"
	"log"

	"github.com/SaputraUta/mini-twitter/services/fanout/internal/model"
	"github.com/SaputraUta/mini-twitter/services/fanout/internal/service"
	"github.com/nats-io/nats.go"
)

type Consumer struct {
	js  nats.JetStreamContext
	svc *service.Fanout
}

func New(js nats.JetStreamContext, svc *service.Fanout) *Consumer {
	return &Consumer{js: js, svc: svc}
}

func (c *Consumer) Start() error {
	_, err := c.js.Subscribe("tweets.created", c.handle, nats.Durable("FANOUT"), nats.ManualAck())
	return err
}

func (c *Consumer) handle(msg *nats.Msg) {
	var ev model.TweetEvent
	if err := json.Unmarshal(msg.Data, &ev); err != nil {
		log.Printf("bad event, dropping %v", err)
		msg.Ack()
		return
	}
	if err := c.svc.Handle(ev); err != nil {
		log.Printf("fanout failed tweet %d: %v", ev.ID, err)
		return
	}
	log.Printf("fanned out tweet %d (user %d)", ev.ID, ev.UserID)
	msg.Ack()
}
