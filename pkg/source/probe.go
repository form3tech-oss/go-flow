package source

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type ProbeSource struct {
	output chan stream.Element
	ctx    context.Context
}

func (t *ProbeSource) AlsoTo(sink stream.Sink) stream.Source {
	panic("implement me")
}

func (t *ProbeSource) DivertTo(sink stream.Sink, when stream.Predicate) stream.Source {
	panic("implement me")
}

func (t *ProbeSource) Via(flow stream.Flow) stream.Source {
	t.output = flow.Input()
	return flow.WireSourceToFlow(t)
}

func (t *ProbeSource) To(sink stream.Sink) stream.Runnable {
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
		t.output <- stream.Value(item)
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
