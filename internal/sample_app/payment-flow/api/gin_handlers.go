package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/internalmodels"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/storage"
	"github.com/form3tech-oss/go-flow/pkg/flow"
	"github.com/form3tech-oss/go-flow/pkg/sink"
	"github.com/form3tech-oss/go-flow/pkg/source"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
)

func handlePayment(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		source.SingleOfHttpRequest(c.Request).
			Via(flow.Map(httpRequestToPayment)).
			//Via(ValidatePayment()).
			To(sink.FromCollector(StorePaymentInPostgres(db))).
			Run(c)
	}
}

func StorePaymentInPostgres(db *sqlx.DB) sink.Collector {
	return &postgresCollector{db: db}
}

func httpRequestToPayment(from types.Element) types.Element {
	request, ok := from.Value.(*http.Request)
	if !ok {
		return types.Error(fmt.Errorf("unexpected type"))
	}

	var payment internalmodels.Payment

	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return types.Error(err)
	}
	err = json.Unmarshal(bytes, &payment)
	if err != nil {
		return types.Error(err)
	}
	return types.Value(payment)
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
