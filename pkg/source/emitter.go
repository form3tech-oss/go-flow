package source

import (
	"context"
	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Emitter interface {
	SetOutput(output chan stream.Element)
	Run(ctx context.Context)
}

type emitterSource struct {
	emitter      Emitter
	divertedSink stream.Sink
	whenToDivert stream.Predicate
}

func (e emitterSource) DivertTo(sink stream.Sink, when stream.Predicate) stream.Source {
	panic("implement me")
}

func (e emitterSource) Via(operation stream.Flow) stream.Source {
	e.emitter.SetOutput(operation.Input())
	return operation.WireSourceToFlow(e)
}

func (e emitterSource) To(sink stream.Sink) stream.Runnable {
	e.emitter.SetOutput(sink.Input())
	return sink.WireSourceToSink(e)
}

func (e emitterSource) Run(ctx context.Context) {
	e.emitter.Run(ctx)
}

func FromEmitter(emitter Emitter) stream.Source {
	return emitterSource{emitter: emitter}
}
