package corel

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var corel interface{} = time.Now().Add(-time.Microsecond * (time.Duration(rand.Intn(1000)))).String()

//CoRelationId correlationData
type CoRelationId struct {
	RequestID string `json:"requestID" binding:"required"`
	SessionID string `json:"sessionID" binding:"required"`
	Hop       int    `json:"hop"`
	isset     bool
	mu        sync.Mutex
}

func (corelid *CoRelationId) OnceMust() {
	if !corelid.isset {
		corelid.mu.Lock()
		if !corelid.isset {
			corelid.isset = true
			corelid.Hop += 1
			if len(corelid.RequestID) == 0 {
				corelid.RequestID = xid.New().String()
			}
			if len(corelid.SessionID) == 0 {
				corelid.SessionID = "null-" + corelid.RequestID
			}
		}
		corelid.mu.Unlock()
	}
}

//GetCorelationId ...
func GetCorelationId(ctx context.Context) (corelid *CoRelationId, err error) {
	var ok bool
	if corelid, ok = ctx.Value(corel).(*CoRelationId); !ok {
		err = fmt.Errorf("invalid corelationId")
	}
	return
}

//GetCtxWithCorelID ...
//Will be of no use once the consumer stops copying the context.
func GetCtxWithCorelID(ctx context.Context, corelid *CoRelationId) context.Context {
	return context.WithValue(ctx, corel, corelid)
}

//GinSetCoRelID ...
func GinSetCoRelID(c *gin.Context, corelid *CoRelationId) {
	c.Set(corel.(string), corelid)
}

func AttachCorelToHttp(corelid *CoRelationId, req *http.Request) {
	req.Header.Set("session_id", corelid.SessionID)
	req.Header.Set("request_id", corelid.RequestID)
	req.Header.Set("hop", fmt.Sprintf("%d", corelid.Hop))
}

func AttachCorelToHttpFromCtx(ctx context.Context, req *http.Request) {
	if corelid, err := GetCorelationId(ctx); err == nil {
		AttachCorelToHttp(corelid, req)
	}
}
