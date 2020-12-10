package source

import (
	"context"
	"testing"
	"time"

	"github.com/form3tech-oss/go-flow/pkg/option"
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

func TestRange_CompletesLongRange(t *testing.T) {
	// Arrange
	defer goleak.VerifyNone(t)
	probe := sink.Probe(t)
	sourceUnderTest := Range(1, 1000000, option.BufferedChannel(100000)).To(probe)

	// Act
	sourceUnderTest.Run(context.Background())

	// Assert
	probe.Request(1000000, 60*time.Second)
	probe.ExpectComplete()
}
