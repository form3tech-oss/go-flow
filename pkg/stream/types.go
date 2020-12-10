package stream

import "context"

type Source interface {
	Output() chan Element
	Via(flow Flow) Source
	To(sink Sink) Runnable
	Runnable
}

type Sink interface {
	SetSource(source Source) Runnable
	Runnable
}

type Flow interface {
	SetSource(source Source) Source
	Source
}

type Runnable interface {
	Run(ctx context.Context)
}

// https://blog.yoshuawuyts.com/rust-streams/
