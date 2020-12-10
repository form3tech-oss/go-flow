package source

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Iterator interface {
	HasNext() bool
	GetNext() stream.Element
}

type iteratorSource struct {
	iterator  Iterator
	output    chan stream.Element
	diversion stream.Sink
	divert    stream.Predicate
}

func (i *iteratorSource) Via(flow stream.Flow) stream.Source {
	i.output = flow.Input()
	return flow.WireSourceToFlow(i)
}

func (i *iteratorSource) To(sink stream.Sink) stream.Runnable {
	i.output = sink.Input()
	return sink.WireSourceToSink(i)
}

func (i *iteratorSource) DivertTo(sink stream.Sink, when stream.Predicate) stream.Source {
	i.diversion = sink
	i.divert = when
	return i
}

func (i *iteratorSource) Run(ctx context.Context) {
	go func() {
		defer i.closeOutputs()
		for i.iterator.HasNext() {
			element := i.iterator.GetNext()
			select {
			case <-ctx.Done():
				return
			default:
				if i.divert(element) {
					i.diversion.Input() <- element
				} else {
					i.output <- element
				}
			}
		}
	}()
}

func (i *iteratorSource) closeOutputs() {
	close(i.output)
	if i.diversion != nil {
		close(i.diversion.Input())
	}
}

func FromIterator(iterator Iterator) stream.Source {
	return &iteratorSource{
		iterator: iterator,
		divert: func(element stream.Element) bool {
			return false
		},
	}
}
