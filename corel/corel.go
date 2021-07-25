package corel

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	"gitlab.com/dotpe/mindbenders/utils/lib/base62"
)

var corel interface{} = time.Now().Add(-time.Microsecond * (time.Duration(rand.Intn(1000)))).String()

//CoRelationId correlationData
type CoRelationId struct {
	RequestID string `json:"requestID" header:"request_id"`
	SessionID string `json:"sessionID" header:"session_id"`
	Hop       int    `json:"hop" header:"hop"`
	Auth      string `header:"Authorization"`
	IsHTTP    bool

	OriginHost string
	OriginApp  string
	JWT        string
	User       *dotJWTinfo

	isset bool
	mu    sync.Mutex
	enc   string
}
type dotJWTinfo struct {
	TenantID  int    `json:"TenantId"`
	StoreID   int    `json:"StoreId"`
	SessionID string `json:"sessionID" header:"session_id"`
	Exp       int    `json:"exp"`
}

func (djwt dotJWTinfo) Valid() error {
	return nil
}
func (corelid *CoRelationId) loadAuth() error {
	if len(corelid.Auth) > 0 && corelid.Auth != "unknownToken" {
		corelid.JWT = strings.Replace(corelid.Auth, "Bearer ", "", 1)
		parts := strings.Split(corelid.JWT, ".")
		raw, err := jwt.DecodeSegment(parts[1])
		if err != nil {
			return err
		}
		var jwtinfo dotJWTinfo
		err = json.Unmarshal(raw, &jwtinfo)
		if err == nil && jwtinfo.TenantID > 0 && jwtinfo.Exp > 0 {
			corelid.User = &jwtinfo
			if len(jwtinfo.SessionID) > 0 {
				corelid.SessionID = fmt.Sprintf("%d:%s", jwtinfo.TenantID, jwtinfo.SessionID)
			} else {
				corelid.SessionID = fmt.Sprintf("%d:%s", jwtinfo.TenantID, base62.Encode(int64(jwtinfo.Exp)))
			}
		}
	}
	return nil
}

func encCorelToBase64(corelid *CoRelationId) string {
	if len(corelid.enc) == 0 {
		corelid.mu.Lock()
		if len(corelid.enc) == 0 {
			raw, _ := json.Marshal(corelid)
			corelid.enc = base64.StdEncoding.EncodeToString(raw)
		}
		corelid.mu.Unlock()
	}
	return corelid.enc
}

func decodeBase64ToCorel(raw string, corel *CoRelationId) error {
	if rawbyte, err := base64.StdEncoding.DecodeString(raw); err != nil {
		return err
	} else {
		return json.Unmarshal(rawbyte, &corel)
	}
}

func (corelid *CoRelationId) OnceMust(c *gin.Context, app string) {
	if !corelid.isset {
		corelid.mu.Lock()
		if !corelid.isset {
			rawcorel := c.Request.Header.Get("corel")
			if len(rawcorel) > 0 {
				decodeBase64ToCorel(rawcorel, corelid)
			}
			if len(corelid.RequestID) == 0 {
				corelid.loadAuth()
				corelid.IsHTTP = true
				corelid.RequestID = xid.New().String()
				corelid.OriginApp = app
				corelid.OriginHost, _ = os.Hostname()
			}
			if len(corelid.SessionID) == 0 {
				corelid.SessionID = "null-" + corelid.RequestID
			}
			corelid.isset = true
			corelid.Hop += 1
		}
		corelid.mu.Unlock()
	}
}

//GetCorelationId ...
func GetCorelationId(ctx context.Context) (corelid *CoRelationId, err error) {
	var ok bool
	if corelid, ok = ctx.Value(corel).(*CoRelationId); !ok {
		err = fmt.Errorf("invalid corelationId")
	}
	return
}

//GetCtxWithCorelID ...
//Will be of no use once the consumer stops copying the context.
func GetCtxWithCorelID(ctx context.Context, corelid *CoRelationId) context.Context {
	return context.WithValue(ctx, corel, corelid)
}

//GinSetCoRelID ...
func GinSetCoRelID(c *gin.Context, corelid *CoRelationId) {
	c.Set(corel.(string), corelid)
}

func AttachCorelToHttp(corelid *CoRelationId, req *http.Request) {
	req.Header.Set("session_id", corelid.SessionID)
	req.Header.Set("request_id", corelid.RequestID)
	req.Header.Set("hop", fmt.Sprintf("%d", corelid.Hop))
	req.Header.Set("corel", encCorelToBase64(corelid))
}

func AttachCorelToHttpFromCtx(ctx context.Context, req *http.Request) {
	if corelid, err := GetCorelationId(ctx); err == nil {
		AttachCorelToHttp(corelid, req)
	}
}
