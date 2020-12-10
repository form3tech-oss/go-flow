package source

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"
)

func Range(start int, end int, options ...option.Option) types.Source {
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

func (r *rangeIterator) GetNext(ctx context.Context) types.Element {
	if r.HasNext(ctx) {
		element := types.Value(r.current)
		r.current++
		return element
	}
	return types.Error(fmt.Errorf("end of stream"))
}
