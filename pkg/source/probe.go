package source

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/types"
)

type ProbeSource struct {
	output chan types.Element
	ctx    context.Context
}

func (t *ProbeSource) AlsoTo(sink types.Sink) types.Source {
	panic("implement me")
}

func (t *ProbeSource) DivertTo(sink types.Sink, when types.Predicate) types.Source {
	panic("implement me")
}

func (t *ProbeSource) Via(flow types.Flow) types.Source {
	t.output = flow.Input()
	return flow.WireSourceToFlow(t)
}

func (t *ProbeSource) To(sink types.Sink) types.Runnable {
	t.output = sink.Input()
	return sink.WireSourceToSink(t)
}

func (t *ProbeSource) Run(ctx context.Context) {
	t.ctx = ctx
}

func (t *ProbeSource) SendAndComplete(items ...interface{}) {
	t.Send(items...)
	t.Complete()
}

func (t *ProbeSource) Send(items ...interface{}) {
	for _, item := range items {
		t.output <- types.Value(item)
	}
}

func (t *ProbeSource) Complete() {
	close(t.output)
}

func Probe() *ProbeSource {
	return &ProbeSource{
		ctx: context.Background(),
	}
}
