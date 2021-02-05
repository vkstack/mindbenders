package id

import (
	"math/rand"
	"sync"

	"github.com/bwmarrin/snowflake"
)

const (
//Miliseconnds for "2019-07-23"
//Never Fuck-up whith this number	CUSTOMEPOCH int64 = 1612137600000
)

var node *snowflake.Node
var once sync.Once

func SetNode(nodeID int64) (err error) {
	snowflake.Epoch = CUSTOMEPOCH
	node, err = snowflake.NewNode(nodeID)
	return
}

func defaultInit() {
	SetNode(rand.Int63n(1024))
}

//GetUID this method returns snowflake nuique id
func GetUID() int64 {
	if node == nil {
		once.Do(defaultInit)
	}
	return node.Generate().Int64()
}
