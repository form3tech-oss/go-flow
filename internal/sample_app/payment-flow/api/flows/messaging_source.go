package flows

import (
	"github.com/form3tech/go-messaging/messaging"
)

type Source struct {
}

func NewMessagingSource() *Source {
	return &Source{}
}

func (f Source) Send(destination string, message messaging.Message) error {
	panic("implement me - send to a flows sync.")
}
