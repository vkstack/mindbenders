package corel

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"gitlab.com/dotpe/mindbenders/errors"

	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var corel interface{} = time.Now().Add(-time.Microsecond * (time.Duration(rand.Intn(1000)))).String()

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

func (corelid *CoRelationId) loadAuth() error {
	if len(corelid.Auth) > 0 && corelid.Auth != "unknownToken" {
		parts := strings.Split(corelid.Auth, ".")
		if len(parts) < 2 {
			return errors.New("invalid auth provided")
		}
		raw, err := jwt.DecodeSegment(parts[1])
		if err != nil {
			return errors.WrapMessage(err, "JWT decoding failed")
		}
		// corelid.User = new(dotJWTinfo)
		return json.Unmarshal(raw, &corelid.JWT)
	}
	return nil
}

func (corelid *CoRelationId) encCorelToBase64() {
	raw, _ := json.Marshal(corelid)
	corelid.enc = base64.StdEncoding.EncodeToString(raw)
}

func decodeBase64ToCorel(raw string, corel *CoRelationId) error {
	if rawbyte, err := base64.StdEncoding.DecodeString(raw); err != nil {
		return errors.WrapMessage(err, "base64 to corel struct decoding failed")
	} else {
		err := json.Unmarshal(rawbyte, &corel)
		return errors.WrapMessage(err, "raw base64 to corel struct unmarshalling failed")
	}
}

func (corelid *CoRelationId) OnceMust(c context.Context) {
	corelid.once.Do(func() {
		if gc, ok := c.(*gin.Context); ok {
			rawcorel := gc.Request.Header.Get("corel")
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
		if len(corelid.SessionID) == 0 {
			corelid.SessionID = "null-" + corelid.RequestID
		}
		corelid.encCorelToBase64()
	})
}

// GetCorelationId ...
func GetCorelationId(ctx context.Context) (corelid *CoRelationId, err error) {
	var ok bool
	if corelid, ok = ctx.Value(corel).(*CoRelationId); !ok {
		err = fmt.Errorf("invalid corelationId")
	}
	return
}

// GetCtxWithCorelID ...
// Will be of no use once the consumer stops copying the context.
func GetCtxWithCorelID(ctx context.Context, corelid *CoRelationId) context.Context {
	return context.WithValue(ctx, corel, corelid)
}

// GinSetCoRelID ...
func GinSetCoRelID(c *gin.Context, corelid *CoRelationId) {
	c.Set(corel.(string), corelid)
}

func AttachCorelToHttp(corelid *CoRelationId, req *http.Request) {
	req.Header.Set("session_id", corelid.SessionID)
	req.Header.Set("request_id", corelid.RequestID)
	req.Header.Set("corel", corelid.enc)
}

func AttachCorelToHttpFromCtx(ctx context.Context, req *http.Request) {
	if corelid, err := GetCorelationId(ctx); err == nil {
		AttachCorelToHttp(corelid, req)
	}
}

func NewCorelCtx(sessionId string) context.Context {
	return NewCorelCtxFromCtx(context.Background(), sessionId)
}

func NewCorelCtxFromCtx(ctx context.Context, sessionId string) context.Context {
	corelId := &CoRelationId{SessionID: sessionId}
	corelId.OnceMust(ctx)
	ctx = context.WithValue(ctx, corel, corelId)
	return ctx
}
