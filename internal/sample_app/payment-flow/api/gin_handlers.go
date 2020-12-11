package api

import (
	"context"
	"fmt"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/internalmodels"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/storage"
	"github.com/form3tech-oss/go-flow/pkg/flow"
	"github.com/form3tech-oss/go-flow/pkg/sink"
	"github.com/form3tech-oss/go-flow/pkg/http"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func handlePayment(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		http.Source(c).
			Via(flow.Map(ContextToPaymentRequest())).
			Via(flow.Map(PaymentIsValid())).
			Via(flow.Map(PaymentPersisted())).
			Via(flow.Map(PaymentPersistedToResponse())).
			To(http.Sink(c)).
			Run(c)
	}
}

func ContextToPaymentRequest() flow.Mapper {
	return func(from types.Element) types.Element {
		request, ok := from.Value.(*gin.Context)
		if !ok {
			return types.Error(fmt.Errorf("unexpected type"))
		}

		var payment internalmodels.Payment

		err := request.BindJSON(&payment)
		if err != nil {
			return types.Error(err)
		}
		return types.Value(payment)
	}
}

func PaymentIsValid() flow.Mapper {
	return func(from types.Element) types.Element {
		return from
	}
}

func PaymentPersisted() flow.Mapper {
	return func(from types.Element) types.Element {
		return from
	}
}

func PaymentPersistedToResponse() flow.Mapper {
	return func(from types.Element) types.Element {
		return from
	}
}


func StorePaymentInPostgres(db *sqlx.DB) sink.Collector {
	return &postgresCollector{db: db}
}


type postgresCollector struct {
	db *sqlx.DB
}

func (c *postgresCollector) Collect(ctx context.Context, element types.Element) {

	// TODO - failing here because element has an error and no value.

	w := storage.GetPaymentWriter(c.db)
	err := w.Create(&ctx, element.Value.(*internalmodels.Payment))
	if err != nil {
		panic(err)
	}
}
