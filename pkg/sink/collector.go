package sink

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/option"

	"github.com/form3tech-oss/go-flow/pkg/types"
)

type collectorSink struct {
	source    types.Source
	collector Collector
	input     chan types.Element
}

func (s *collectorSink) Input() chan types.Element {
	return s.input
}

type Collector interface {
	Collect(ctx context.Context, element types.Element)
}

func (s *collectorSink) WireSourceToSink(source types.Source) types.Runnable {
	s.source = source
	return s
}

func (s *collectorSink) Run(ctx context.Context) {
	s.source.Run(ctx)
	go func() {
		for s.input != nil {
			select {
			case <-ctx.Done():
				return
			case element, ok := <-s.input:
				if !ok {
					return
				}
				s.collector.Collect(ctx, element)
			default:
			}
		}
	}()
}

func FromCollector(collector Collector, options ...option.Option) types.Sink {
	return &collectorSink{
		collector: collector,
		input:     option.CreateChannel(options...),
	}
}
