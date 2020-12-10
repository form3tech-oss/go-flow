package stream

import "context"

type Predicate func(element Element) bool

type Source interface {
	Via(flow Flow) Source
	To(sink Sink) Runnable
	DivertTo(sink Sink, when Predicate) Source
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
	Source
}

type Runnable interface {
	Run(ctx context.Context)
}
