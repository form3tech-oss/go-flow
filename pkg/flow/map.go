package flow

import (
	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type Mapper func(from stream.Element) stream.Element

type mappingOperator struct {
	mapper Mapper
}

func (m *mappingOperator) apply(element stream.Element) (stream.Element, bool) {
	mappedElement := m.mapper(element)
	return mappedElement, true
}

func Map(mapper Mapper, options ...option.Option) stream.Flow {
	return FromOperator(&mappingOperator{
		mapper: mapper,
	})
}
