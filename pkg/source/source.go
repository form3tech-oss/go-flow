package source

import (
	"context"
	"github.com/form3tech-oss/go-flow/pkg/types"
)

type source struct {
	iterator  types.Iterator
	output    chan types.Element
	diversion types.Sink
	divert    types.Predicate
	alsoTo    types.Sink
}

func (s *source) AlsoTo(sink types.Sink) types.Source {
	s.alsoTo = sink
	return s
}

func (s *source) Via(flow types.Flow) types.Source {
	s.output = flow.Input()
	return flow.WireSourceToFlow(s)
}

func (s *source) To(sink types.Sink) types.Runnable {
	s.output = sink.Input()
	return sink.WireSourceToSink(s)
}

func (s *source) DivertTo(sink types.Sink, when types.Predicate) types.Source {
	s.diversion = sink
	s.divert = when
	return s
}

func (s *source) Run(ctx context.Context) {
	go func() {
		defer s.closeOutputs()
		for s.iterator.HasNext(ctx) {
			element := s.iterator.GetNext(ctx)
			select {
			case <-ctx.Done():
				return
			default:
				if s.divert(element) {
					s.diversion.Input() <- element
				} else {
					s.output <- element
				}
			}
		}
	}()
}

func (s *source) runAttachedStages(ctx context.Context) {
	if s.diversion != nil {
		s.diversion.Run(ctx)
	}
	if s.alsoTo != nil {
		s.alsoTo.Run(ctx)
	}
}

func (s *source) closeOutputs() {
	close(s.output)
	if s.diversion != nil {
		close(s.diversion.Input())
	}
}

