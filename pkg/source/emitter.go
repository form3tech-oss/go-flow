package source

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Emitter interface {
	Output() chan stream.Element
	Run(ctx context.Context)
}

type emitterSource struct {
	emitter Emitter
}

func (e emitterSource) Output() chan stream.Element {
	return e.emitter.Output()
}

func (e emitterSource) Via(operation stream.Flow) stream.Source {
	return operation.SetSource(e)
}

func (e emitterSource) To(sink stream.Sink) stream.Runnable {
	return sink.SetSource(e)
}

func (e emitterSource) Run(ctx context.Context) {
	e.emitter.Run(ctx)
}

func FromEmitter(emitter Emitter) stream.Source {
	return emitterSource{emitter: emitter}
}
