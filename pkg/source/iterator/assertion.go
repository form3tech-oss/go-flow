package iterator

import (
	"context"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func expectToMatch(t *testing.T, iterator types.Iterator, expected ... interface{}) {

	ctx := context.Background()

	index :=0

	for iterator.HasNext(ctx) {
		actualItem := iterator.GetNext(ctx)
		assert.True(t, index <= len(expected), "expected the iterator to have the same number of items as expected" )
		expectedItem := types.Value(expected[index])
		assert.Equal(t, expectedItem.Value, actualItem.Value, "expected to iterate the same value")
		index ++
	}

}


