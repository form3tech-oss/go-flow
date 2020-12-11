package source

import (
	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/source/iterator"
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
	return FromIterator(iterator.FromEmitter(emitter))
}

func Range(start int, end int, options ...option.Option) types.Source {
	return FromIterator(&rangeIterator{
		start:   start,
		end:     end,
		current: start,
	})
}

func SingleOfInt(value ...int) types.Source {
	return FromIterator(iterator.OfInts(value...))
}

func SingleOfString(value ...string) types.Source {
	return FromIterator(iterator.OfStrings(value...))
}

func OfInts(value ...int) types.Source {
	return FromIterator(iterator.OfInts(value...))
}
