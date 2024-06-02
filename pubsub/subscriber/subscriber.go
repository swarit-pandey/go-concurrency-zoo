package subscriber

import (
	"fmt"

	"github.com/swarit-pandey/go-concurrency-zoo/pubsub/producer"
)

type ConcreteSubscriber struct {
	ID string
}

func (cs *ConcreteSubscriber) Notify(msg producer.Message) {
	fmt.Printf("Subscriber %s received message on topic %s: %v\n", cs.ID, msg.Topic, msg.Payload)
}
