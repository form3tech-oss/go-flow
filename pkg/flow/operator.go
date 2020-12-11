package flow

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"
)

type Operator interface {
	apply(element types.Element) (types.Element, bool)
}

type operatorFlow struct {
	source    types.Source
	operator  Operator
	input     chan types.Element
	output    chan types.Element
	diversion types.Sink
	divert    types.Predicate
	alsoTo    types.Sink
}

func (o *operatorFlow) AlsoTo(sink types.Sink) types.Source {
	o.alsoTo = sink
	return o
}

func (o *operatorFlow) Input() chan types.Element {
	return o.input
}

func (o *operatorFlow) DivertTo(sink types.Sink, when types.Predicate) types.Source {
	o.diversion = sink
	o.divert = when
	return o
}

func (o *operatorFlow) WireSourceToFlow(source types.Source) types.Source {
	o.source = source
	return o
}

func (o *operatorFlow) Via(flow types.Flow) types.Source {
	o.output = flow.Input()
	return flow.WireSourceToFlow(o)
}

func (o *operatorFlow) To(sink types.Sink) types.Runnable {
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

func FromOperator(operator Operator, options ...option.Option) types.Flow {
	return &operatorFlow{
		operator: operator,
		input:    option.CreateChannel(options...),
		divert: func(element types.Element) bool {
			return false
		},
	}
}
