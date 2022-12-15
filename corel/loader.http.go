package corel

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
* inside header it will write `Corel`="base64 of json"
 */
func HttpCorelLoader(ctx context.Context, header http.Header) {
	if corelid, err := GetCorelationId(ctx); err == nil {
		header.Set("request_id", corelid.RequestId)
		header.Set("session_id", corelid.SessionId)
		header.Set(string(CtxCorelLocator), corelid.Child().Enc())
	}
}

func HttpCorelUnLoader(ctx context.Context, header http.Header) context.Context {
	corelid := DecodeCorelationId(header.Get(string(CtxCorelLocator)))
	corelid1, _ := corel(ctx)
	if !strings.HasPrefix(corelid1.SessionId, "null") {
		corelid = corelid1
	}
	if len(corelid.SessionId) == 0 {
		return ctx
	}
	if gc, ok := ctx.(*gin.Context); ok {
		gc.Set(string(CtxCorelLocator), corelid)
		return gc
	}
	return context.WithValue(ctx, CtxCorelLocator, corelid)
}

/*
	(t1,t2,t3)
	* http loaders
	[]func(ctx context.Context, header *http.Header)
	* http unloaders
	[]func(ctx context.Context, header *http.Header) context.Context


	* Kafka loaders
	[]func(ctx context.Context, header *[]kafka.Header)
	* Kafka unloaders
	[]func(ctx context.Context, header *[]kafka.Header) context.Context

*/
