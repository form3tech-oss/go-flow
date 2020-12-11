package source

import (
	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"
)

func FromIterator(iterator types.Iterator) types.Source {
	return &source{
		iterator: iterator,
		divert: func(element types.Element) bool {
			return false
		},
	}
}

func FromEmitter(emitter types.Emitter) types.Source {
	return FromIterator(emitterIterator{
		hasStarted: false,
		emitter:    emitter,
	})
}


func Range(start int, end int, options ...option.Option) types.Source {
	return FromIterator(&rangeIterator{
		start:   start,
		end:     end,
		current: start,
	})
}