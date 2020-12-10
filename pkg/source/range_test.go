package source

import (
	"context"
	"testing"
	"time"

	"github.com/form3tech-oss/go-flow/pkg/option"
	"github.com/form3tech-oss/go-flow/pkg/types"

	"github.com/form3tech-oss/go-flow/pkg/sink"
	"go.uber.org/goleak"
)

func TestRange_ReadsExpectedAndCompletes(t *testing.T) {
	// Arrange
	defer goleak.VerifyNone(t)
	probe := sink.Probe(t)
	sourceUnderTest := Range(1, 10).To(probe)

	// Act
	sourceUnderTest.Run(context.Background())

	// Assert
	probe.Request(10, 1*time.Second)
	probe.Expect(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	probe.ExpectComplete()
}

func TestRange_CanBeDivertedAndCompletes(t *testing.T) {
	// Arrange
	defer goleak.VerifyNone(t)
	minProbe := sink.Probe(t)
	maxProbe := sink.Probe(t)
	sourceUnderTest := Range(1, 10).DivertTo(maxProbe, func(element types.Element) bool {
		return element.Value.(int) > 5
	}).To(minProbe)

	// Act
	sourceUnderTest.Run(context.Background())

	// Assert
	minProbe.Request(5, 5*time.Second)
	minProbe.Expect(1, 2, 3, 4, 5)
	maxProbe.Request(5, 5*time.Second)
	maxProbe.Expect(6, 7, 8, 9, 10)
	minProbe.ExpectComplete()
	maxProbe.ExpectComplete()

}

func TestRange_CompletesLongRange(t *testing.T) {
	// Arrange
	defer goleak.VerifyNone(t)
	probe := sink.Probe(t, option.BufferedChannel(100000))
	sourceUnderTest := Range(1, 1000000).To(probe)

	// Act
	sourceUnderTest.Run(context.Background())

	// Assert
	probe.Request(1000000, 60*time.Second)
	probe.ExpectComplete()
}
