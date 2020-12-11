package types

import "context"

type Predicate func(element Element) bool

type Source interface {
	Via(flow Flow) Source
	To(sink Sink) Runnable
	DivertTo(sink Sink, when Predicate) Source
	AlsoTo(sink Sink) Source
	Runnable
}

type Sink interface {
	Input() chan Element
	WireSourceToSink(source Source) Runnable
	Runnable
}

type Flow interface {
	Input() chan Element
	WireSourceToFlow(source Source) Source
	DivertTo(sink Sink, when Predicate) Source
	AlsoTo(sink Sink) Source
	Source
}

type Runnable interface {
	Run(ctx context.Context)
}

type Emitter interface {
	Output() chan Element
	Run(ctx context.Context)
}

type Iterator interface {
	HasNext(ctx context.Context) bool
	GetNext(ctx context.Context) Element
}
