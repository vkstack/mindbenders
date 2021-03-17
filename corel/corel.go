package corel

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var corel interface{} = time.Now().Add(-time.Microsecond * (time.Duration(rand.Intn(1000)))).String()

//CoRelationId correlationData
type CoRelationId struct {
	RequestID string `json:"requestID"`
	SessionID string `json:"sessionID"`
}

//GetCorelationId ...
func GetCorelationId(ctx context.Context) (corelid CoRelationId, err error) {
	var ok bool
	if corelid, ok = ctx.Value(corel).(CoRelationId); !ok {
		err = fmt.Errorf("invalid corelationId")
	}
	return
}

//GetCtxWithCorelID ...
//Will be of no use once the consumer stops copying the context.
func GetCtxWithCorelID(ctx context.Context, requestID, sessionID string) context.Context {
	return context.WithValue(ctx, corel, CoRelationId{
		RequestID: requestID,
		SessionID: sessionID,
	})
}

//GinSetCoRelID ...
func GinSetCoRelID(c *gin.Context, corelid CoRelationId) {
	c.Set(corel.(string), corelid)
}
