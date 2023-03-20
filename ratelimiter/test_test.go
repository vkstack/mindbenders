package ratelimiter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

func TestLimiter(t *testing.T) {
	rl := redis_rate.NewLimiter(
		redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{"localhost:6379"},
		}),
	)
	t1 := time.NewTicker(time.Second)
	t2 := time.NewTicker(2 * time.Second)
	lim := redis_rate.Limit{
		Period: time.Second * 15,
		Rate:   3,
		Burst:  3,
	}
	curr := time.Now()
	// stopat := 2 * time.Minute
	for {
		select {
		case x := <-t1.C:
			fmt.Println(int(x.Sub(curr)/time.Second), Eval(context.Background(), rl, lim, "test"))
		case y := <-t2.C:
			if y.Sub(curr) >= 2*time.Minute {
				t1.Stop()
				t2.Stop()
			}
			break
		}
		if t1.C == nil && t2.C == nil {
			break
		}
	}
}
