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

var (
	ctxcorelLocator corelstr = "corelid-local"
	corelHeaderKey           = "corel"
)

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

func (corelid *CoRelationId) encCorelToBase64() {
	raw, _ := json.Marshal(corelid)
	corelid.enc = base64.StdEncoding.EncodeToString(raw)
}

func decodeBase64ToCorel(raw string, corel *CoRelationId) error {
	rawbyte, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return errors.WrapMessage(err, "base64 to corel struct decoding failed")
	}
	if err := json.Unmarshal(rawbyte, &corel); err != nil {
		return errors.WrapMessage(err, "raw base64 to corel struct unmarshalling failed")
	}
	corel.enc = raw
	return nil
}

func corel(ctx context.Context) (*CoRelationId, error) {
	if corelid, ok := ctx.Value(ctxcorelLocator).(*CoRelationId); ok {
		return corelid, nil
	}
	if c, ok := ctx.(*gin.Context); ok {
		if v, ok := c.Get(string(ctxcorelLocator)); ok {
			if corelid, ok := v.(*CoRelationId); ok {
				return corelid, nil
			}
		}
		var corelid = new(CoRelationId)
		c.ShouldBindHeader(&corelid)
		corelid.init(c)
		if len(corelid.RequestID) > 0 {
			c.Set(string(ctxcorelLocator), corelid)
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
	req.Header.Set(corelHeaderKey, corelid.child().enc)
}

func AttachCorelToHttpFromCtx(ctx context.Context, req *http.Request) {
	if corelid, err := corel(ctx); err == nil {
		AttachCorelToHttp(corelid, req)
	}
}
