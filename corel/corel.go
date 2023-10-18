// The code defines a Go package called "corel" that provides functionality for generating and managing
// correlation IDs.
package corel

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/xid"
)

// CoRelationId correlationData
type CoRelationId struct {
	SessionId,
	RequestId,
	AppRequestId,
	RequestSource,
	enc string
}

func (corelid *CoRelationId) Enc() string { return corelid.enc }

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
		SessionId:     "null-" + tmp,
		RequestId:     tmp,
		AppRequestId:  tmp,
		RequestSource: "",
	}
	switch x := len(ids); {
	case x > 3:
		corelid.RequestSource = ids[3]
		fallthrough
	case x > 2:
		corelid.AppRequestId = ids[2]
		fallthrough
	case x > 1:
		corelid.RequestId = ids[1]
		fallthrough
	case x > 0:
		corelid.SessionId = ids[0]
	}
	corelid.enc = EncodeCorel(&corelid)
	return &corelid
}

func NewCorelIdFromHttp(header http.Header) *CoRelationId {
	if sessid := header.Get("session_id"); len(sessid) > 0 {
		return NewCorelId(sessid)
	}
	if auth := header.Get("Authorization"); len(auth) > 0 {
		if sessId := getJWTSession(auth); len(sessId) > 0 {
			return NewCorelId(sessId)
		}
	}
	if rawcorel := header.Get(string(CtxCorelLocator)); len(rawcorel) > 0 {
		var corelid *CoRelationId
		if err := DecodeCorel(rawcorel, &corelid); err == nil {
			corelid.enc = rawcorel
			return corelid
		}
	}
	return NewCorelId()
}
