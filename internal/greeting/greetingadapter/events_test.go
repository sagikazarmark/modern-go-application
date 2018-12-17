package greetingadapter

import "github.com/ThreeDotsLabs/watermill/message"

type publisherStub struct {
	topic    string
	messages []*message.Message
}

func (p *publisherStub) Publish(topic string, messages ...*message.Message) error {
	p.topic = topic
	p.messages = messages

	return nil
}

func (*publisherStub) Close() error {
	return nil
}
