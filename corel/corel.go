package corel

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/errors"
)

// CoRelationId correlationData
type CoRelationId struct {
	RequestID string `json:"requestID" header:"request_id"`
	SessionID string `json:"sessionID" header:"session_id"`
	Auth      string `header:"Authorization"`
	JWT       *jwtinfo

	AppRequestId,
	RequestSource string

	enc string

	once sync.Once
}

type jwtinfo struct {
	SessionID string `json:"sessionID" header:"session_id" validate:"required"`
}

func (corelid *CoRelationId) init(c context.Context) {
	corelid.once.Do(func() {
		if gc, ok := c.(*gin.Context); ok {
			gc.ShouldBindHeader(&corelid)
			rawcorel := gc.Request.Header.Get(corelHeaderKey)
			if len(rawcorel) > 0 {
				if err := decodeBase64ToCorel(rawcorel, corelid); err == nil {
					return
				}
			}
		}
		if len(corelid.RequestId()) == 0 {
			corelid.loadAuth()
			corelid.RequestID = xid.New().String()
		}
		if len(corelid.AppRequestId) == 0 {
			corelid.AppRequestId = xid.New().String()
		}
		if corelid.JWT != nil && corelid.JWT.SessionID != "" {
			corelid.SessionID = corelid.JWT.SessionID
		} else if len(corelid.SessionID) == 0 {
			corelid.SessionID = "null-" + corelid.RequestID
		}
		corelid.encCorelToBase64()
	})
}

// todo: do we need this
func NewCoRelationId(sessionId string) *CoRelationId {
	corelid := &CoRelationId{RequestID: xid.New().String(), SessionID: sessionId}
	if sessionId == "" {
		corelid.SessionID = corelid.RequestID
	}
	return corelid
}

func (corelid *CoRelationId) child() *CoRelationId {
	ch := CoRelationId(*corelid)
	ch.AppRequestId = xid.New().String()
	ch.RequestSource = os.Getenv("APP") + ":" + corelid.AppRequestId
	ch.encCorelToBase64()
	return &ch
}

func (jwt *jwtinfo) UnmarshalJSON(raw []byte) error {
	var x struct {
		SessionID string `json:"sessionID"`
	}
	if err := json.Unmarshal(raw, &x); err != nil {
		return err
	}
	if x.SessionID == "" {
		return errors.New("invalid sessionId")
	}
	*jwt = jwtinfo(x)
	return nil
}

// This is used when the default context is used to define a new corel
func NewCorelCtx(sessionId string) context.Context {
	return NewCorelCtxFromCtx(context.Background(), sessionId)
}

// This is used to define a new corel on the context
func NewCorelCtxFromCtx(ctx context.Context, sessionId string) context.Context {
	corelId := &CoRelationId{SessionID: sessionId}
	corelId.init(ctx)
	ctx = context.WithValue(ctx, ctxcorelLocator, corelId)
	return ctx
}

// Any task consumer will call this method only
func NewCorelCtxFromRequest(ctx context.Context, sessionId, requestId string) context.Context {
	corelId := &CoRelationId{SessionID: sessionId, RequestID: requestId, AppRequestId: xid.New().String()}
	corelId.init(ctx)
	ctx = context.WithValue(ctx, ctxcorelLocator, corelId)
	return ctx
}

func (corelid *CoRelationId) Logrus(f logrus.Fields) {
	f["sessionId"] = corelid.SessionID
	f["requestId"] = corelid.RequestID
	f["appRequestId"] = corelid.AppRequestId
	if len(corelid.RequestSource) != 0 {
		f["requestSource"] = corelid.RequestSource
	}
}

func (corelid *CoRelationId) SessionId() string {
	return corelid.SessionID
}
func (corelid *CoRelationId) RequestId() string {
	return corelid.RequestID
}
