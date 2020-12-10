package flow

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Operator interface {
	SetInput(input chan stream.Element)
	Output() chan stream.Element
	Run(ctx context.Context)
}

type operatorFlow struct {
	source   stream.Source
	operator Operator
}

func (o *operatorFlow) SetSource(source stream.Source) stream.Source {
	o.source = source
	o.operator.SetInput(source.Output())
	return o
}

func (o *operatorFlow) Output() chan stream.Element {
	return o.operator.Output()
}

func (o *operatorFlow) Via(flow stream.Flow) stream.Source {
	return flow.SetSource(o)
}

func (o *operatorFlow) To(sink stream.Sink) stream.Runnable {
	return sink.SetSource(o)
}

func (o *operatorFlow) Run(ctx context.Context) {
	o.source.Run(ctx)
	o.operator.Run(ctx)
}

func FromOperator(operator Operator) stream.Flow {
	return &operatorFlow{
		operator: operator,
	}
}
