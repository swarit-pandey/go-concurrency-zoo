package producer

import "sync"

type Publisher struct {
	mu          sync.Mutex
	subscribers map[string][]Subscriber
}

func New() *Publisher {
	return &Publisher{
		subscribers: make(map[string][]Subscriber),
	}
}

func (p *Publisher) Subscribe(topic string, subscriber Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subscribers[topic] = append(p.subscribers[topic], subscriber)
}

func (p *Publisher) Publish(topic string, payload any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if subs, found := p.subscribers[topic]; found {
		msg := &Message{Topic: topic, Payload: payload}
		for _, sub := range subs {
			sub.Notify(*msg)
		}
	}
}
