package events

import (
	"encoding/json"

	"github.com/SaputraUta/mini-twitter/services/posts/internal/model"
	"github.com/nats-io/nats.go"
)

const tweetCreatedSubject = "tweets.created"

type NatsPublisher struct {
	js nats.JetStreamContext
}

func NewNatsPublisher(nc *nats.Conn) (*NatsPublisher, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	if _, err := js.StreamInfo("TWEETS"); err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "TWEETS",
			Subjects: []string{tweetCreatedSubject},
		})
		if err != nil {
			return nil, err
		}
	}

	return &NatsPublisher{js: js}, nil
}

func (p *NatsPublisher) PublishTweetCreated(t model.Tweet) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	_, err = p.js.Publish(tweetCreatedSubject, data)
	return err
}
