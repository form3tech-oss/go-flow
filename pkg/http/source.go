package http

import (
	"github.com/form3tech-oss/go-flow/pkg/source"
	"github.com/form3tech-oss/go-flow/pkg/source/iterator"
	"github.com/form3tech-oss/go-flow/pkg/types"
	"github.com/gin-gonic/gin"
)

func OfGinContexts(value ...*gin.Context) types.Iterator {
	return  iterator.OfAny(ginContextToAny(value...))
}

func ginContextToAny(in ...*gin.Context) []interface{} {
	out := make([]interface{}, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}


func Source(context *gin.Context) types.Source {
	return  source.FromIterator(OfGinContexts(context))
}


