package flow

import (
	"context"
	"testing"
	"time"

	"github.com/form3tech-oss/go-flow/pkg/stream"
	"github.com/form3tech-oss/go-flow/pkg/stream/sink"
	"github.com/form3tech-oss/go-flow/pkg/stream/source"
	"go.uber.org/goleak"
)

func TestMappingOperator_Run(t *testing.T) {
	defer goleak.VerifyNone(t)

	sourceProbe := source.Probe()
	sinkProbe := sink.Probe(t)

	flowUnderTest := Map(func(element stream.Element) stream.Element {
		return stream.Value(element.Value.(int) * 2)
	})

	sourceProbe.Via(flowUnderTest).To(sinkProbe).Run(context.Background())

	go func() {
		sourceProbe.SendNext(1)
		sourceProbe.SendNext(6)
		sourceProbe.SendNext(4)
		sourceProbe.Complete()
	}()

	sinkProbe.Request(3, 1*time.Second)
	sinkProbe.Expect(2, 12, 8)

}
