package source

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Emitter interface {
	Output() chan stream.Element
	Run(ctx context.Context)
}

type emitterIterator struct {
	hasStarted bool
	emitter Emitter
	current stream.Element
}

func (e emitterIterator) HasNext(ctx context.Context) bool {
	if !e.hasStarted {
		e.emitter.Run(ctx)
	}
	element, ok := <-e.emitter.Output()
	if ok {
		e.current = element
	}
	return ok
}

func (e emitterIterator) GetNext(ctx context.Context) stream.Element {
	return e.current
}

func FromEmitter(emitter Emitter) stream.Source {
	return FromIterator(emitterIterator{})
}
