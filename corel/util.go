package corel

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"gitlab.com/dotpe/mindbenders/errors"

	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type corelstr string

const CtxCorelLocator corelstr = "corel"

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
		return json.Unmarshal(raw, &corelid.JWT)
	}
	return nil
}

func corel(ctx context.Context) (*CoRelationId, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}
	if corelid, ok := ctx.Value(CtxCorelLocator).(*CoRelationId); ok {
		return corelid, nil
	}
	if c, ok := ctx.(*gin.Context); ok {
		if v, ok := c.Get(string(CtxCorelLocator)); ok {
			if corelid, ok := v.(*CoRelationId); ok {
				return corelid, nil
			}
		}
		var corelid = new(CoRelationId)
		corelid.init(c)
		if len(corelid.RequestId) > 0 {
			c.Set(string(CtxCorelLocator), corelid)
			// this will help the other middleware to copy request headers
			c.Header(string(CtxCorelLocator), corelid.enc)
		}
		return corelid, nil
	}
	return nil, errors.New("invalid/missing corelationId")
}

// GetCorelationId ...
func GetCorelationId(ctx context.Context) (corelid *CoRelationId, err error) {
	return corel(ctx)
}

func AttachCorelToHttp(corelid *CoRelationId, req *http.Request) {
	req.Header.Set("request_id", corelid.RequestId)
	req.Header.Set("session_id", corelid.SessionId)
	req.Header.Set(string(CtxCorelLocator), corelid.Child().enc)
}

func AttachCorelToHttpFromCtx(ctx context.Context, req *http.Request) {
	if corelid, err := corel(ctx); err == nil {
		AttachCorelToHttp(corelid, req)
	}
}

func EncodeCorel(corelId *CoRelationId) string {
	raw, _ := json.Marshal(corelId)
	return base64.StdEncoding.EncodeToString(raw)
}

func DecodeCorel(str string, dst interface{}) error {
	rawbyte, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.WrapMessage(err, "base64 to corel struct decoding failed")
	}
	if err := json.Unmarshal(rawbyte, &dst); err != nil {
		return errors.WrapMessage(err, "raw base64 to corel struct unmarshalling failed")
	}
	return nil
}

func DecodeCorelationId(encoded string) *CoRelationId {
	var corel CoRelationId
	if err := DecodeCorel(encoded, &corel); err != nil {
		corelid, _ := GetCorelationId(NewCorelCtx(""))
		return corelid
	}
	corel.enc = encoded
	corel.once.Do(func() {})
	return &corel
}
