package flow

import (
	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"
)

type Mapper func(from types.Element) types.Element

type mappingOperator struct {
	mapper Mapper
}

func (m *mappingOperator) apply(element types.Element) (types.Element, bool) {
	mappedElement := m.mapper(element)
	return mappedElement, true
}

func Map(mapper Mapper, options ...option.Option) types.Flow {
	return FromOperator(&mappingOperator{
		mapper: mapper,
	})
}
