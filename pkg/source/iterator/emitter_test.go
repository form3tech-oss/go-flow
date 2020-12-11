package iterator

import (
	"context"
	"testing"

	"github.com/form3tech-oss/go-flow/pkg/types"
	"go.uber.org/goleak"
)

type testEmitter struct {
	output   chan types.Element
	iterator types.Iterator
	t        *testing.T
}

func createIntEmitterIterator(t *testing.T, values ...int) types.Iterator {
	return FromEmitter(&testEmitter{
		output:   make(chan types.Element),
		iterator: OfInts(values...),
		t:        t,
	})
}

func (e *testEmitter) Output() chan types.Element {
	return e.output
}

func (e *testEmitter) Run(ctx context.Context) {
	go func() {
		index := 0
		for e.iterator.HasNext(ctx) {
			index++
			e.t.Logf("next is %v", index)
			e.output <- e.iterator.GetNext(ctx)
		}
		close(e.output)
	}()
}

func TestEmitterIterator_EmitsAllAsExpected(t *testing.T) {
	defer goleak.VerifyNone(t)
	iterator := createIntEmitterIterator(t, 1, 2, 3)
	expectToMatch(t, iterator, 1, 2, 3)
}
