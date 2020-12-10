package sink

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type collectorSink struct {
	source    stream.Source
	collector Collector
}

type Collector interface {
	SetInput(input chan stream.Element)
	Run(ctx context.Context)
}

func (s *collectorSink) SetSource(source stream.Source) stream.Runnable {
	s.source = source
	s.collector.SetInput(s.source.Output())
	return s
}

func (s *collectorSink) Run(ctx context.Context) {
	s.source.Run(ctx)
	s.collector.Run(ctx)
}

func FromCollector(collector Collector) stream.Sink {
	return &collectorSink{
		collector: collector,
	}
}
