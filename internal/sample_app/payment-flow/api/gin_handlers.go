package api

import (
	"context"
	"fmt"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/internalmodels"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/storage"
	"github.com/form3tech-oss/go-flow/pkg/flow"
	"github.com/form3tech-oss/go-flow/pkg/http"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func handlePayment(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		completed := make(chan struct{})

		http.Source(c).
			Via(flow.Map(ContextToPaymentRequest())).
			Via(flow.Map(PaymentIsValid())).
			Via(flow.Map(PaymentPersisted(context.Background(), db))).
			Via(flow.Map(PaymentPersistedToResponse())).
			To(http.Sink(c, completed)).
			Run(c)

		select {
			case <- completed  :
			case <-c.Done() :
		}
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
		payment := from.Value.(*internalmodels.Payment)
		if payment.Record.Status == "duplicate" {
			return types.Error( fmt.Errorf("duplicate payment , %v", payment.ID))
		}
		return from
	}
}


func PaymentPersisted(ctx context.Context, db *sqlx.DB) flow.Mapper {
	return func(from types.Element) types.Element {

		if from.Error != nil {
			return from
		}

		w := storage.GetPaymentWriter(db)
		err := w.Create(&ctx, from.Value.(*internalmodels.Payment))
		if err != nil {
			return types.Error(err)
		}
		return from
	}
}

func PaymentPersistedToResponse() flow.Mapper {
	return func(from types.Element) types.Element {


		if from.Error != nil {
		  return types.Value(http.Response{
				StatusCode: 400,
				Body:       from.Error.Error(),
			}		)
		}

		types.Value(http.Response{
			StatusCode: 204,
			Body:       nil,
		})

		return from
	}
}


