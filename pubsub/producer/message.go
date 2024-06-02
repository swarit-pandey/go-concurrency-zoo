package producer

type Message struct {
	Topic   string
	Payload any
}

type Subscriber interface {
	Notify(msg Message)
}
