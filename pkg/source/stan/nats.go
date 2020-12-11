package stan

import (
    "context"
	"github.com/form3tech-oss/go-flow/pkg/source"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/nats-io/stan.go"
	"log"
)

type natsEmitter struct {
	conn                stan.Conn
	subject             string
	group               string
	output              chan types.Element
	subscriptionOptions []stan.SubscriptionOption
}


func (r *natsEmitter) Output() chan types.Element {
	return r.output
}

func (r *natsEmitter) Run(ctx context.Context) {
	go func() {
		var sub stan.Subscription
		var err error

		defer func() {
			close(r.output)
			if sub != nil {
				sub.Close()
			}
		}()

		sub, err = r.conn.QueueSubscribe(r.subject, r.group, func(msg *stan.Msg) {
			select {
			case <-ctx.Done():
				log.Printf("Returning ...:  %v ", ctx.Err())
				return
			case  r.output <- types.Value(msg.Data):
				log.Printf("Received: %v", msg.Data )
				err := msg.Ack()
				if err != nil {
					log.Printf("ack failed: %v", err)
				}
			}
		}, r.subscriptionOptions...)

		if err != nil {
			log.Printf("failed consuming from nats: %v", err)
			return
		}

		select {
		case <- ctx.Done():
			log.Printf("Returning context at end...:  %v ", ctx.Err())
		}

	}()
}

func Source(conn stan.Conn, group string, subject string, subscriptionOptions []stan.SubscriptionOption, options ... option.Option) types.Source {
	return source.FromEmitter(&natsEmitter{
		conn:                conn,
		group:               group,
		subject:             subject,
		subscriptionOptions: subscriptionOptions,
		output:              option.CreateChannel(options...),
	})
}
