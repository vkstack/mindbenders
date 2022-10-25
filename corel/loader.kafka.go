package corel

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

// unloading of corel info from ctx and loading it in kafka header
func KafkaCorelLoader(ctx context.Context, headers *[]kafka.Header) {
	if corelid, err := GetCorelationId(ctx); err == nil {
		*headers = append(*headers, kafka.Header{
			Key:   string(CtxCorelLocator),
			Value: []byte(corelid.Child().Enc()),
		})
	}
}

// This will be used to load corel from kafka message's header to context
// unloading of corel info from header and loading it in ctx
func KafkaCorelUnLoader(ctx context.Context, headers *[]kafka.Header) context.Context {
	for _, headr := range *headers {
		if headr.Key == string(CtxCorelLocator) {
			corelid := DecodeCorelationId(string(headr.Value)).Sibling()
			if gc, ok := ctx.(*gin.Context); ok {
				gc.Set(string(CtxCorelLocator), corelid)
				return gc
			}
			return context.WithValue(ctx, CtxCorelLocator, corelid)
		}
	}
	return ctx
}
