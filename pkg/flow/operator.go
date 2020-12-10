package flow

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Operator interface {
	apply(element stream.Element) (stream.Element, bool)
}

type operatorFlow struct {
	source    stream.Source
	operator  Operator
	input     chan stream.Element
	output    chan stream.Element
	diversion stream.Sink
	divert    stream.Predicate
	alsoTo stream.Sink
}


func (o *operatorFlow) AlsoTo(sink stream.Sink) stream.Source {
	o.alsoTo = sink
	return o
}

func (o *operatorFlow) Input() chan stream.Element {
	return o.input
}

func (o *operatorFlow) DivertTo(sink stream.Sink, when stream.Predicate) stream.Source {
	o.diversion = sink
	o.divert = when
	return o
}

func (o *operatorFlow) WireSourceToFlow(source stream.Source) stream.Source {
	o.source = source
	return o
}

func (o *operatorFlow) Via(flow stream.Flow) stream.Source {
	o.output = flow.Input()
	return flow.WireSourceToFlow(o)
}

func (o *operatorFlow) To(sink stream.Sink) stream.Runnable {
	o.output = sink.Input()
	return sink.WireSourceToSink(o)
}

func (o *operatorFlow) Run(ctx context.Context) {
	o.runAttachedStages(ctx)
	go func() {
		defer o.closeOutputs()
		for o.input != nil {
			select {
			case <-ctx.Done():
				return
			case element, ok := <-o.input:
				if !ok {
					return
				}
				resultElement, ok := o.operator.apply(element)
				if ok {
					if o.divert(resultElement) {
						o.diversion.Input() <- resultElement
					} else {
						o.output <- resultElement
						if o.alsoTo != nil {
							o.alsoTo.Input() <- resultElement
						}
					}
				}
			default:

			}
		}
	}()
}

func (o *operatorFlow) runAttachedStages(ctx context.Context) {
	if o.source != nil {
		o.source.Run(ctx)
	}
	if o.diversion != nil {
		o.diversion.Run(ctx)
	}
	if o.alsoTo != nil {
		o.alsoTo.Run(ctx)
	}
}

func (o *operatorFlow) closeOutputs() {
	close(o.output)
	if o.diversion != nil {
		close(o.diversion.Input())
	}
	if o.alsoTo != nil {
		close(o.alsoTo.Input())
	}
}

func FromOperator(operator Operator, options ...option.Option) stream.Flow {
	return &operatorFlow{
		operator: operator,
		input:    option.CreateChannel(options...),
		divert: func(element stream.Element) bool {
			return false
		},
	}
}
