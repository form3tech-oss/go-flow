package sink

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type probeSink struct {
	source stream.Source
	items  []interface{}
	t      *testing.T
	ctx    context.Context
}

func (p *probeSink) SetSource(source stream.Source) stream.Runnable {
	p.source = source
	return p
}

func (p *probeSink) Run(ctx context.Context) {
	p.source.Run(ctx)
	p.ctx = ctx
}

func (p *probeSink) Request(number int, timeout time.Duration) {
	var wg sync.WaitGroup
	wg.Add(1)
	p.items = nil
	go func() {
		defer wg.Done()
		ti := time.After(timeout)
		for p.source.Output() != nil {
			select {
			case <-p.ctx.Done():
				return
			case element, ok := <-p.source.Output():
				if !ok {
					return
				}
				p.items = append(p.items, element.Value)
				if len(p.items) == number {
					return
				}
			case <-ti:
				p.t.Errorf("Timed out receiving")
				return
			default:
			}
		}
	}()
	wg.Wait()
}

func (p *probeSink) Expect(expected ...interface{}) {
	if !testEq(p.items, expected) {
		p.t.Logf("Expected next to contain the same expected items \n expected: %v\n actual %v", expected, p.items)
		p.t.Fail()
	}
}

func (p *probeSink) ExpectComplete() {
	_, ok := <-p.source.Output()
	if ok {
		p.t.Error("Expected the source out put channel to be closed.")
	}
}

func Probe(t *testing.T) *probeSink {
	return &probeSink{
		source: nil,
		t:      t,
		ctx:    context.Background(),
	}
}

func testEq(a, b []interface{}) bool {

	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
