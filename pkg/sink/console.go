package sink

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"
)

type consoleCollector struct {
}

func (c *consoleCollector) Collect(ctx context.Context, element types.Element) {
	fmt.Println(element)
}

func Console(options ...option.Option) types.Sink {
	return FromCollector(&consoleCollector{}, options...)
}
