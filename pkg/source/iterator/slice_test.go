package iterator

import (
	"testing"
)

func TestSliceIterator_ToCompleteCorrectly(t *testing.T) {
	expectToMatch(t, OfInts(1, 2, 3, 4), 1, 2, 3, 4)
	expectToMatch(t, OfStrings("a", "B", "C", "e"), "a", "B", "C", "e")
}
