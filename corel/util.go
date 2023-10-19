package corel

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

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

var encoder = base64.StdEncoding

const encsep = "\t"

func EncodeCorel(corelid *CoRelationId) string {
	if len(corelid.enc) == 0 {
		corelid.enc = fmt.Sprintf("%s%s%s%s%s%s%s", corelid.SessionId, encsep, corelid.RequestId, encsep, corelid.AppRequestId, encsep, corelid.RequestSource)
		corelid.enc = encoder.EncodeToString([]byte(corelid.enc))
	}
	return corelid.enc
}

func DecodeCorel(encoded []byte) (corelid *CoRelationId, err error) {
	var decoded []byte = make([]byte, encoder.DecodedLen(len(encoded)))
	n, err := encoder.Decode(decoded, encoded)
	if err != nil {
		return nil, err
	}
	decoded = decoded[:n]
	if parts := bytes.Split(decoded, []byte(encsep)); len(parts) >= 4 {
		corelid = new(CoRelationId)
		corelid.enc = string(encoded)
		corelid.SessionId, corelid.RequestId = string(parts[0]), string(parts[1])
		corelid.AppRequestId, corelid.RequestSource = string(parts[2]), string(parts[3])
		return corelid, nil
	}
	return nil, errors.New("invalid encoded corel")
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
