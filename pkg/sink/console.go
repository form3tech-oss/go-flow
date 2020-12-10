package sink

import (
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/stream"
)

type consoleCollector struct {
}

func (c *consoleCollector) Collect(element stream.Element) {
	fmt.Println(element)
}

func Console(options ...option.Option) stream.Sink {
	return FromCollector(&consoleCollector{}, options...)
}
