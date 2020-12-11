package stan

import (
	"context"
	"github.com/form3tech-oss/go-flow/pkg/flow"
	"github.com/form3tech-oss/go-flow/pkg/sink"
	"github.com/form3tech-oss/go-flow/pkg/source"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"go.uber.org/goleak"
	"strconv"
	"testing"
	"time"
)

type testStanCollector struct {
	t         *testing.T
	conn      stan.Conn
	subject   string
	complete  chan struct{}
	lastValue int
}

func (t *testStanCollector) Collect(ctx context.Context, element types.Element) {
	if element.Error != nil {
		return
	}

	err := t.conn.Publish("stan_subject", []byte(strconv.Itoa(element.Value.(int))))
	if err != nil{
		t.t.Fatal(err)
	}

	if element.Value.(int) == t.lastValue {
		t.complete <- struct{}{}
	}
}

func newPublishSink(t *testing.T, conn stan.Conn, subject string, lastValue int, complete chan struct{}) types.Sink {
	return sink.FromCollector(& testStanCollector{
		t:       t,
		conn:    conn,
		subject: subject,
		lastValue : lastValue,
		complete : complete,
	})
}


func getStanConnection(t * testing.T) stan.Conn {
	sc, err := stan.Connect("nats-streaming", "stan_streaming_publisher", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		t.Fatal(err)
	}
	return sc
}


func TestStream(t *testing.T) {

	defer goleak.VerifyNone(t)
	probe := sink.Probe(t)

	sc := getStanConnection(t)
	defer sc.Close()

	subscriptionOptions := getSubscriptionOptions()

	subject := uuid.New().String()

	complete := make(chan struct{})
	// Publishing
	source.OfInts(1,2,3,4,5,6,7,8,9,10).To(newPublishSink(t, sc, subject, 10, complete)).Run(context.Background())
	 <- complete


	sourceUnderTest := Source(sc, "stan_streaming_test", subject, subscriptionOptions).Via(flow.Map(mapBytesToInts(t)))
	sourceUnderTest.To(probe)

	// Act
	sourceUnderTest.Run(context.Background())

	// Assert
	probe.Request(10, 20*time.Second)

	probe.Expect(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	//probe.ExpectComplete()
}

func mapBytesToInts(t *testing.T) flow.Mapper {
	return func(from types.Element) types.Element {
		value, err := strconv.Atoi(string(from.Value.([]byte)))
		if err != nil {
			t.Fatal(err)
		}
		return types.Value(value)
	}
}

func getSubscriptionOptions() []stan.SubscriptionOption {
	subscriptionOptions := []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.DurableName("stan_streaming_test"),
		stan.StartWithLastReceived(),
		stan.MaxInflight(1),
	}
	return subscriptionOptions
}
