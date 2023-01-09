package corel

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

// unloading of corel info from ctx and loading it in kafka header
func AmqpLoader(ctx context.Context, headers amqp.Table) {
	if corelid, err := GetCorelationId(ctx); err == nil {
		headers[string(CtxCorelLocator)] = corelid.Child().Enc()
	}
}

// This will be used to load corel from kafka message's header to context
// unloading of corel info from header and loading it in ctx
func AmqpUnloader(ctx context.Context, headers amqp.Table) context.Context {
	if h, ok := headers[string(CtxCorelLocator)]; ok {
		if raw, ok := h.(string); ok {
			corelid := DecodeCorelationId(raw).Sibling()
			if gc, ok := ctx.(*gin.Context); ok {
				gc.Set(string(CtxCorelLocator), corelid)
				return gc
			}
			return context.WithValue(ctx, CtxCorelLocator, corelid)
		}
	}
	return ctx
}
