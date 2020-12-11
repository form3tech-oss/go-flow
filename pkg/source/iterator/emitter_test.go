package iterator

import (
	"context"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"go.uber.org/goleak"
	"testing"
)

type testEmitter struct {
	output chan types.Element
	iterator types.Iterator
}

func createIntEmitterIterator(values ... int) types.Iterator {
	return FromEmitter(&testEmitter{
		output:   make(chan types.Element),
		iterator: OfInts(values ...),
	})
}

func (e *testEmitter) Output() chan types.Element {
	return e.output
}

func (e *testEmitter) Run(ctx context.Context) {
	go func() {
		for e.iterator.HasNext(ctx) {
			e.output <- e.iterator.GetNext(ctx)
		}
		close(e.output)
	}()
}

func TestEmitterIterator_EmitsAllAsExpected(t *testing.T) {
	defer goleak.VerifyNone(t)
	iterator := createIntEmitterIterator(1,2,3)
	expectToMatch(t, iterator, 1,2,3)
}

