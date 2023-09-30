package corel

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// unloading of corel info from ctx and loading it in kafka header
func KafkaCorelLoader(ctx context.Context, headers []kafka.Header) []kafka.Header {
	cheader := kafka.Header{
		Key:   string(CtxCorelLocator),
		Value: []byte(GetCorelationId(ctx).Child().Enc()),
	}
	if len(headers) == 0 {
		return []kafka.Header{cheader}
	}
	for i := range headers {
		if headers[i].Key == string(CtxCorelLocator) {
			headers[i].Value = cheader.Value
			return headers
		}
	}
	headers = append(headers, cheader)
	return headers
}

// This will be used to load corel from kafka message's header to context
// unloading of corel info from header and loading it in ctx
func KafkaCorelUnLoader(ctx context.Context, headers *[]kafka.Header) context.Context {
	for _, headr := range *headers {
		if headr.Key == string(CtxCorelLocator) {
			corelid := DecodeCorelationId(string(headr.Value)).Sibling()
			return context.WithValue(ctx, CtxCorelLocator, corelid)
		}
	}
	return ctx
}
