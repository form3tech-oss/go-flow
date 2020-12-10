package flow

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Mapper func(from stream.Element) stream.Element

type mappingOperator struct {
	output chan stream.Element
	input  chan stream.Element
	mapper Mapper
}

func (m *mappingOperator) SetInput(input chan stream.Element) {
	m.input = input
}

func (m *mappingOperator) Output() chan stream.Element {
	return m.output
}

func (m *mappingOperator) Run(ctx context.Context) {
	go func() {
		for m.input != nil {
			select {
			case <-ctx.Done():
				return
			case element, ok := <-m.input:
				if !ok {
					return
				}
				m.output <- m.mapper(element)
			default:
			}
		}
	}()
}

func Map(mapper Mapper, options ...option.Option) stream.Flow {
	// channel should be built from options.
	return FromOperator(&mappingOperator{
		output: option.CreateChannel(options...),
		mapper: mapper,
	})
}
