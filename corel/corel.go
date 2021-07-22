package corel

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

	isset bool
	mu    sync.Mutex
}

type dotJWTinfo struct {
	// Username      string `json:"username"`
	TenantID int `json:"TenantId"`
	// StoreID       int    `json:"StoreId"`
	// FeatureRoleID string `json:"FeatureRoleId"`
	// ExpiryTime    string `json:"ExpiryTime"`
	// IssueTime     string `json:"IssueTime"`
	// UserType      string `json:"UserType"`
	SessionID string `json:"sessionID" header:"session_id"`
	Exp       int    `json:"exp"`
	JWT
	// Iss           string `json:"iss"`
}

func (corelid *CoRelationId) OnceMust() {
	if !corelid.isset {
		corelid.mu.Lock()
		if !corelid.isset {
			corelid.isset = true
			corelid.Hop += 1
			if len(corelid.RequestID) == 0 {
				corelid.RequestID = xid.New().String()
			}
			if len(corelid.Auth) > 0 && corelid.Auth != "unknownToken" {
				corelid.Auth = strings.Replace(corelid.Auth, "Bearer ", "", 1)
				_, strs, _ := new(jwt.Parser).ParseUnverified(corelid.Auth, jwt.MapClaims{})
				var jwtinfo dotJWTinfo
				json.Unmarshal([]byte(strs[1]), &jwtinfo)
				if jwtinfo.TenantID > 0 && jwtinfo.Exp > 0 {
					corelid.SessionID = fmt.Sprintf("%d:%s", jwtinfo.TenantID, base62.Encode(int64(jwtinfo.Exp)))
				}
			}
			if len(corelid.SessionID) == 0 {
				corelid.SessionID = "null-" + corelid.RequestID
			}
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
}

func AttachCorelToHttpFromCtx(ctx context.Context, req *http.Request) {
	if corelid, err := GetCorelationId(ctx); err == nil {
		AttachCorelToHttp(corelid, req)
	}
}
