package id

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	UNUSEDBITS   = 1
	EPOCHBITS    = 41
	NODEIDBITS   = 10
	SEQUENCEBITS = 12

	MAXNODEID  = 1023
	MAXSEQUECE = 4095

	// miliseconnds for "2019-07-23"
	CUSTOMEPOCH int64 = 1563840000000
)

var counter int64 = 0

// GetUIDFromNodeCounter sh
func GetUIDFromNodeCounter(nodeID, counter int64) int64 {
	nowmili := time.Now().UnixNano() / 1e6
	customMili := nowmili - CUSTOMEPOCH
	nodeID &= MAXNODEID
	var id int64 = (customMili) << (NODEIDBITS + SEQUENCEBITS)
	id |= (int64(nodeID << SEQUENCEBITS))
	id |= int64(counter)
	return id
}

func GetUIDFromNode(nodeID int64) int64 {
	seq := atomic.AddInt64(&counter, 1) & MAXSEQUECE
	return GetUIDFromNodeCounter(nodeID, seq)
}

func GetUIDFromCounter(counter int64) int64 {
	return GetUIDFromNodeCounter(rand.Int63n(1024), counter)
}

func GetUID() int64 {
	seq := atomic.AddInt64(&counter, 1) & MAXSEQUECE
	return GetUIDFromNodeCounter(rand.Int63n(1024), seq)
}