package iterator

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/types"
)

type sliceIterator struct {
	items        []interface{}
	currentIndex int
}

func (s *sliceIterator) HasNext(ctx context.Context) bool {
	return s.currentIndex < len(s.items)
}

func (s *sliceIterator) GetNext(ctx context.Context) types.Element {
	defer func() { s.currentIndex++ }()
	if s.HasNext(ctx) {
		return types.Value(s.items[s.currentIndex])
	}
	return types.Error(fmt.Errorf("reached end of slice"))
}

func ofAny(value []interface{}) types.Iterator {
	return &sliceIterator{
		items:        value,
		currentIndex: 0,
	}
}
