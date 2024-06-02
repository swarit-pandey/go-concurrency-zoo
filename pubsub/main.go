package main

import (
	"github.com/swarit-pandey/go-concurrency-zoo/pubsub/producer"
	"github.com/swarit-pandey/go-concurrency-zoo/pubsub/subscriber"
)

func main() {
	publisher := producer.New()

	sub1 := &subscriber.ConcreteSubscriber{ID: "sub1"}
	sub2 := &subscriber.ConcreteSubscriber{ID: "sub2"}

	publisher.Subscribe("topic1", sub1)
	publisher.Subscribe("topic1", sub2)

	publisher.Publish("topic1", "hello world")
}
