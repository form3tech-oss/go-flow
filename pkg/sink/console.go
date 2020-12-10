package sink

import (
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"
)

type consoleCollector struct {
}

func (c *consoleCollector) Collect(element types.Element) {
	fmt.Println(element)
}

func Console(options ...option.Option) types.Sink {
	return FromCollector(&consoleCollector{}, options...)
}
