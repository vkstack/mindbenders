package corel

import (
	"context"

	"github.com/streadway/amqp"
)

// unloading of corel info from ctx and loading it in kafka header
func AmqpLoader(ctx context.Context, headers amqp.Table) amqp.Table {
	corelid := GetCorelationId(ctx)
	headers[string(CtxCorelLocator)] = EncodeCorel(corelid.Child())
	return headers
}

// This will be used to load corel from kafka message's header to context
// unloading of corel info from header and loading it in ctx
func AmqpUnloader(ctx context.Context, headers amqp.Table) context.Context {
	if h, ok := headers[string(CtxCorelLocator)]; ok {
		var corelid *CoRelationId
		switch v := h.(type) {
		case string:
			corelid, _ = DecodeCorel([]byte(v))
		case []byte:
			corelid, _ = DecodeCorel(v)
		}

		if corelid == nil {
			corelid = NewCorelId()
		}
		return context.WithValue(ctx, CtxCorelLocator, corelid)
	}
	return ctx
}
