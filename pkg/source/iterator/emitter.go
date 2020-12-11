package iterator

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/types"
)

type emitterIterator struct {
	hasStarted bool
	emitter    types.Emitter
	current    types.Element
}

func (e *emitterIterator) HasNext(ctx context.Context) bool {
	if !e.hasStarted {
		e.emitter.Run(ctx)
		e.hasStarted = true
	}
	element, ok := <-e.emitter.Output()
	if ok {
		e.current = element
	}
	return ok
}

func (e *emitterIterator) GetNext(ctx context.Context) types.Element {
	return e.current
}


