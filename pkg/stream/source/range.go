package source

import (
	"context"

	"github.com/form3tech-oss/go-flow/pkg/stream"
	"github.com/form3tech-oss/go-flow/pkg/stream/option"
)

type rangeEmitter struct {
	start  int
	end    int
	output chan stream.Element
}

func (r *rangeEmitter) Output() chan stream.Element {
	return r.output
}

func (r *rangeEmitter) Run(ctx context.Context) {
	go func() {
		for number := r.start; number <= r.end; number++ {
			select {
			case <-ctx.Done():
				return
			case r.output <- stream.Value(number):
			}
		}
		close(r.output)
	}()
}

func Range(start int, end int, options ...option.Option) stream.Source {
	return FromEmitter(&rangeEmitter{
		start:  start,
		end:    end,
		output: option.CreateChannel(options...),
	})
}
