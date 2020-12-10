package sink

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/option"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type collectorSink struct {
	source    stream.Source
	collector Collector
	input     chan stream.Element
}

func (s *collectorSink) Input() chan stream.Element {
	return s.input
}

type Collector interface {
	Collect(element stream.Element)
}

func (s *collectorSink) WireSourceToSink(source stream.Source) stream.Runnable {
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
				s.collector.Collect(element)
			default:
			}
		}
	}()
}

func FromCollector(collector Collector, options ...option.Option) stream.Sink {
	return &collectorSink{
		collector: collector,
		input:     option.CreateChannel(options...),
	}
}
