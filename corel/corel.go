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
	RequestId string `json:"requestId" header:"request_id"`
	SessionId string `json:"sessionId" header:"session_id"`
	Auth      string `header:"Authorization"`
	JWT       *jwtinfo

	AppRequestId,
	RequestSource string

	enc string

	once sync.Once
}

func (corelid *CoRelationId) Enc() string { return corelid.enc }

type jwtinfo struct {
	SessionId string `json:"sessionID" header:"session_id" validate:"required"`
}

func (corelid *CoRelationId) init(c context.Context) {
	corelid.once.Do(func() {
		if gc, ok := c.(*gin.Context); ok {
			gc.ShouldBindHeader(&corelid)
			rawcorel := gc.Request.Header.Get(string(CtxCorelLocator))
			if len(rawcorel) > 0 {
				if err := DecodeCorel(rawcorel, corelid); err == nil {
					return
				}
				corelid.enc = rawcorel
			}
		}
		if len(corelid.RequestId) == 0 {
			corelid.loadAuth()
			corelid.RequestId = xid.New().String()
		}
		if len(corelid.AppRequestId) == 0 {
			corelid.AppRequestId = xid.New().String()
		}
		if corelid.JWT != nil && corelid.JWT.SessionId != "" {
			corelid.SessionId = corelid.JWT.SessionId
		} else if len(corelid.SessionId) == 0 {
			corelid.SessionId = "null-" + corelid.RequestId
		}
		corelid.enc = EncodeCorel(corelid)
	})
}

func (corelid *CoRelationId) Child() *CoRelationId {
	ch := CoRelationId(*corelid)
	ch.AppRequestId = xid.New().String()
	ch.RequestSource = os.Getenv("APP") + ":" + corelid.AppRequestId
	ch.enc = EncodeCorel(&ch)
	return &ch
}

func (corelid *CoRelationId) Sibling() *CoRelationId {
	ch := CoRelationId(*corelid)
	ch.AppRequestId = xid.New().String()
	ch.enc = EncodeCorel(&ch)
	return &ch
}

func (jwt *jwtinfo) UnmarshalJSON(raw []byte) error {
	var x struct {
		SessionId string `json:"sessionID"`
	}
	if err := json.Unmarshal(raw, &x); err != nil {
		return err
	}
	if x.SessionId == "" {
		return errors.New("invalid sessionId")
	}
	*jwt = jwtinfo(x)
	return nil
}

// This is used when the default context is used to define a new corel
func NewCorelCtx(sessionId string) context.Context {
	return NewCorelCtxFromCtx(context.Background(), sessionId)
}

func NewCorelCtxFromCorel(corelid *CoRelationId) context.Context {
	return context.WithValue(context.Background(), CtxCorelLocator, corelid)
}

// This is used to define a new corel on the context
func NewCorelCtxFromCtx(ctx context.Context, sessionId string) context.Context {
	corelId := &CoRelationId{SessionId: sessionId}
	corelId.init(ctx)
	ctx = context.WithValue(ctx, CtxCorelLocator, corelId)
	return ctx
}

// Any task consumer will call this method only
func NewCorelCtxFromRequest(ctx context.Context, sessionId, requestId string) context.Context {
	corelId := &CoRelationId{SessionId: sessionId, RequestId: requestId, AppRequestId: xid.New().String()}
	corelId.init(ctx)
	ctx = context.WithValue(ctx, CtxCorelLocator, corelId)
	return ctx
}

func (corelid *CoRelationId) Logrus(f logrus.Fields) {
	f["sessionId"] = corelid.SessionId
	f["requestId"] = corelid.RequestId
	f["appRequestId"] = corelid.AppRequestId
	if len(corelid.RequestSource) != 0 {
		f["requestSource"] = corelid.RequestSource
	}
}

func (corelid *CoRelationId) GetSessionId() string {
	return corelid.SessionId
}

func (corelid *CoRelationId) GetRequestId() string {
	return corelid.RequestId
}
