package corel

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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

func (corelid *CoRelationId) init(ctx context.Context) {
	corelid.once.Do(func() {
		if gc, ok := ctx.(*gin.Context); ok {
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
	return NewCorelId(
		corelid.SessionId,
		corelid.RequestId,
		xid.New().String(),
		fmt.Sprintf("%s-%s", os.Getenv("APP"), corelid.AppRequestId),
	)
}

func (corelid *CoRelationId) Sibling() *CoRelationId {
	return NewCorelId(
		corelid.SessionId,
		corelid.RequestId,
		xid.New().String(),
	)
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

func (corelid *CoRelationId) GetSessionId() string {
	return corelid.SessionId
}

func (corelid *CoRelationId) GetRequestId() string {
	return corelid.RequestId
}

/*
ids[0] - sessionId
ids[1] - requestId
ids[3] - subRequestId
ids[4] - requestSource
*/
func NewCorelId(ids ...string) *CoRelationId {
	var tmp = xid.New().String()
	var corelid = CoRelationId{
		SessionId:     tmp,
		RequestId:     tmp,
		AppRequestId:  tmp,
		RequestSource: "",
	}
	switch x := len(ids); {
	case x >= 4:
		corelid.RequestSource = ids[3]
		fallthrough
	case x >= 3:
		corelid.AppRequestId = ids[2]
		fallthrough
	case x >= 2:
		corelid.RequestId = ids[1]
		fallthrough
	case x >= 1:
		corelid.SessionId = ids[0]
	}
	corelid.init(context.TODO())
	return &corelid
}
