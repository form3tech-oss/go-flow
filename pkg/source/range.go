package source

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type rangeEmitter struct {
	start  int
	end    int
	output chan stream.Element
}

func (r *rangeEmitter) SetOutput(output chan stream.Element) {
	r.output = output
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
	return FromIterator(&rangeIterator{
		start:   start,
		end:     end,
		current: start,
	})
}

/// Range Iterator

type rangeIterator struct {
	start   int
	current int
	end     int
}

func (r *rangeIterator) HasNext() bool {
	return r.current <= r.end
}

func (r *rangeIterator) GetNext() stream.Element {
	if r.HasNext() {
		element := stream.Value(r.current)
		r.current++
		return element
	}
	return stream.Error(fmt.Errorf("end of stream"))
}
