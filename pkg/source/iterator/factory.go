package iterator

import "github.com/form3tech-oss/go-flow/pkg/types"

func Single(value int) types.Iterator {
	return ofAny(intsToAny(value))
}

func OfStrings(value ...string) types.Iterator {
	return ofAny(stringsToAny(value...))
}

func OfInts(value ...int) types.Iterator {
	return ofAny(intsToAny(value...))
}

func FromEmitter(emitter types.Emitter) types.Iterator {
	return &emitterIterator{
		hasStarted: false,
		emitter:    emitter,
	}
}

func stringsToAny(in ...string) []interface{} {
	out := make([]interface{}, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}

func intsToAny(in ...int) []interface{} {
	out := make([]interface{}, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}
