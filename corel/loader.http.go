package corel

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
* inside header it will write `Corel`="base64 of json"
 */
func HttpCorelLoader(ctx context.Context, header http.Header) http.Header {
	corelid := GetCorelationId(ctx)
	header.Set(string(CtxCorelLocator), EncodeCorel(corelid.Child()))
	return header
}

/*
*The below Unloader will work with gin.Context only.
 */
func HttpCorelUnLoader(ctx context.Context, header http.Header) context.Context {
	corelid := NewCorelIdFromHttp(header)
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
