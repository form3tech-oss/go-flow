package http

import (
	"context"
	"github.com/form3tech-oss/go-flow/pkg/sink"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)


type Response struct {
	StatusCode int
	Body interface{}
}

func Sink (c *gin.Context) types.Sink {
	return sink.FromCollector( & responseCollector{c: c})
}


type responseCollector struct {
	c *gin.Context
}

func (r responseCollector) Collect(ctx context.Context, element types.Element) {
	if element.Error != nil {
		r.c.JSON(http.StatusInternalServerError, nil)
		return
	}
	response := element.Value.(Response)
	r.c.JSON(response.StatusCode, response.Body)
}

