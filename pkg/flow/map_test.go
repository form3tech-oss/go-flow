package flow

import (
	"context"
	"testing"
	"time"

	"github.com/form3tech-oss/go-flow/pkg/sink"
	"github.com/form3tech-oss/go-flow/pkg/source"
	"github.com/form3tech-oss/go-flow/pkg/stream"
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
		sourceProbe.SendAndComplete(1, 6, 4)
	}()

	sinkProbe.Request(3, 1*time.Second)
	sinkProbe.Expect(2, 12, 8)

}

func TestMappingOperator_WithDiversion(t *testing.T) {
	defer goleak.VerifyNone(t)

	sourceProbe := source.Probe()
	terminationProbe := sink.Probe(t)
	divertProbe := sink.Probe(t)

	mapping := func(element stream.Element) stream.Element {
		return stream.Value(element.Value.(int) * 2)
	}

	whenGreaterThan10 := func(element stream.Element) bool {
		return element.Value.(int) > 5
	}

	sourceProbe.Via(Map(mapping)).DivertTo(divertProbe, whenGreaterThan10).To(terminationProbe).Run(context.Background())

	go func() {
		sourceProbe.SendAndComplete(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	}()

	terminationProbe.Request(2, 10*time.Second)
	divertProbe.Request(8, 10*time.Second)
	terminationProbe.Expect(2, 4)
	divertProbe.Expect(6, 8, 10, 12, 14, 16, 18, 20)
	divertProbe.ExpectComplete()
	terminationProbe.ExpectComplete()

}
