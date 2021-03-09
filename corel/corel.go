package corel

import (
	"context"
	"fmt"
	"time"
)

var corel interface{} = time.Now()

//GetCorelKey ...
func GetCorelKey() interface{} {
	return corel
}

//CoRelationId correlationData
type CoRelationId struct {
	RequestID string `json:"requestID"`
	SessionID string `json:"sessionID"`
}

//GetCorelationId ...
func GetCorelationId(ctx context.Context) (corelid CoRelationId, err error) {
	var ok bool
	if corelid, ok = ctx.Value(corel).(CoRelationId); !ok {
		err = fmt.Errorf("invalid corelationId")
	}
	return
}
