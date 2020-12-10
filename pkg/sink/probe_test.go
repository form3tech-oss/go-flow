package sink

import (
	"context"
	"testing"
	"time"

	"github.com/form3tech-oss/go-flow/pkg/source"

	"go.uber.org/goleak"
)

func TestProbe_ReadsExpectedAndCompletes(t *testing.T) {
	defer goleak.VerifyNone(t)

	source := source.Probe()
	sinkUnderTest := Probe(t)

	source.To(sinkUnderTest).Run(context.Background())

	go func() {
		source.SendNext(1)
		source.SendNext(2)
		source.SendNext(3)
		source.SendNext(4)
		source.SendNext(5)
		source.SendNext(6)
		source.Complete()
	}()

	sinkUnderTest.Request(1, 1*time.Second)
	sinkUnderTest.Expect(1)
	sinkUnderTest.Request(5, 1*time.Second)
	sinkUnderTest.Expect(2, 3, 4, 5, 6)
	sinkUnderTest.ExpectComplete()

}
