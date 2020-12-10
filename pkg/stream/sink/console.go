package sink

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/stream"
	"github.com/form3tech-oss/go-flow/pkg/stream/option"
)

type consoleCollector struct {
	input chan stream.Element
}

func (c *consoleCollector) SetInput(input chan stream.Element) {
	c.input = input
}

func (c *consoleCollector) Run(ctx context.Context) {
	go func() {
		for c.input != nil {
			select {
			case <-ctx.Done():
				return
			case element, ok := <-c.input:
				if !ok {
					return
				}
				fmt.Println(element)
			default:
			}
		}
	}()
}

func Console(options ...option.Option) stream.Sink {
	return FromCollector(&consoleCollector{})
}
