package corel

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"gitlab.com/dotpe/mindbenders/errors"
)

// CoRelationId correlationData
type CoRelationId struct {
	RequestID string `json:"requestID" header:"request_id"`
	SessionID string `json:"sessionID" header:"session_id"`
	Auth      string `header:"Authorization"`
	JWT       *jwtinfo

	OriginHost,
	OriginApp,
	enc string

	once sync.Once
}

func (corelid *CoRelationId) init(c context.Context) {
	corelid.once.Do(func() {
		if gc, ok := c.(*gin.Context); ok {
			rawcorel := gc.Request.Header.Get(corelHeaderKey)
			if len(rawcorel) > 0 {
				decodeBase64ToCorel(rawcorel, corelid)
			}
		}
		if len(corelid.RequestID) == 0 {
			corelid.loadAuth()
			corelid.RequestID = xid.New().String()
			corelid.OriginApp = os.Getenv("APP")
			corelid.OriginHost, _ = os.Hostname()
		}
		if corelid.JWT != nil && corelid.JWT.SessionID != "" {
			corelid.SessionID = corelid.JWT.SessionID
		} else if len(corelid.SessionID) == 0 {
			corelid.SessionID = "null-" + corelid.RequestID
		}
		corelid.encCorelToBase64()
	})
}

func NewCoRelationId(sessionId string) *CoRelationId {
	corelid := &CoRelationId{RequestID: xid.New().String(), SessionID: sessionId}
	if sessionId == "" {
		corelid.SessionID = corelid.RequestID
	}
	return corelid
}

type jwtinfo struct {
	SessionID string `json:"sessionID" header:"session_id" validate:"required"`
}

func (corelid *CoRelationId) NewChild() *CoRelationId {
	ch := CoRelationId(*corelid)
	//todo:  define prent - child relationship here
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

func NewCorelCtx(sessionId string) context.Context {
	return NewCorelCtxFromCtx(context.Background(), sessionId)
}

func NewCorelCtxFromCtx(ctx context.Context, sessionId string) context.Context {
	corelId := &CoRelationId{SessionID: sessionId}
	corelId.init(ctx)
	ctx = context.WithValue(ctx, ctxcorelLocator, corelId)
	return ctx
}
