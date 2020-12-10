package source

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/stream"
)

func Range(start int, end int, options ...option.Option) stream.Source {
	return FromIterator(&rangeIterator{
		start:   start,
		end:     end,
		current: start,
	})
}

type rangeIterator struct {
	start   int
	current int
	end     int
}

func (r *rangeIterator) HasNext(ctx context.Context) bool {
	return r.current <= r.end
}

func (r *rangeIterator) GetNext(ctx context.Context) stream.Element {
	if r.HasNext(ctx) {
		element := stream.Value(r.current)
		r.current++
		return element
	}
	return stream.Error(fmt.Errorf("end of stream"))
}
