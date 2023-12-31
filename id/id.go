package id

import (
	"math/rand"
	"sync"

	"github.com/bwmarrin/snowflake"
)

const (
	//Miliseconnds for "2021-01-01"
	//Never Fuck-up whith this number
	CUSTOMEPOCH int64 = 1612137600000
)

var node *snowflake.Node
var once sync.Once

func GetGenerator(nodeID int64) (*snowflake.Node, error) {
	snowflake.Epoch = CUSTOMEPOCH
	return snowflake.NewNode(nodeID)
}

//SetNode create a snowflake node.
func SetNode(nodeID int64) (err error) {
	snowflake.Epoch = CUSTOMEPOCH
	node, err = GetGenerator(nodeID)
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
