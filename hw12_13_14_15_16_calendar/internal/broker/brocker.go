package queue

type Broker interface {
	Send(Message) error
	Receive(Message) error
}

type Message struct {
	Message string
}
