package source

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type ProbeSource struct {
	output chan stream.Element
	ctx    context.Context
}

func (t *ProbeSource) Output() chan stream.Element {
	return t.output
}

func (t *ProbeSource) Via(operation stream.Flow) stream.Source {
	return operation.SetSource(t)
}

func (t *ProbeSource) To(sink stream.Sink) stream.Runnable {
	return sink.SetSource(t)
}

func (t *ProbeSource) Run(ctx context.Context) {
	t.ctx = ctx
}

func (t *ProbeSource) SendNext(item interface{}) {
	t.output <- stream.Value(item)
}

func (t *ProbeSource) Complete() {
	close(t.output)
}

func Probe() *ProbeSource {
	return &ProbeSource{
		output: make(chan stream.Element),
		ctx:    context.Background(),
	}
}
