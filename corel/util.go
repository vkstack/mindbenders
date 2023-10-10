package corel

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"gitlab.com/dotpe/mindbenders/errors"

	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type corelstr string

const CtxCorelLocator corelstr = "corel"

func ReadCorelId(ctx context.Context) (*CoRelationId, error) {
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
	}
	return nil, errors.New("invalid/missing corelationId")
}

// concurrent unsafe
// it adds corelId if not found
func GetCorelationId(ctx context.Context) *CoRelationId {
	corelid, _ := ReadCorelId(ctx)
	if corelid == nil {
		if c, ok := ctx.(*gin.Context); ok {
			corelid = NewCorelIdFromHttp(c.Request.Header)
		} else {
			corelid = NewCorelId()
		}
		if c, ok := ctx.(*gin.Context); ok {
			c.Set(string(CtxCorelLocator), corelid)
		}
	}
	return corelid
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
		return NewCorelId()
	}
	corel.enc = encoded
	return &corel
}

func getJWTSession(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return ""
	}
	raw, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return ""
	}
	var jwtinfo struct {
		SessionId string `json:"sessionID" header:"session_id" validate:"required"`
	}
	if err := json.Unmarshal(raw, &jwtinfo); err != nil {
		return ""
	}
	return jwtinfo.SessionId
}
