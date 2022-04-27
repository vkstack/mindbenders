package ratelimiter

// func TestLimiter(t *testing.T) {
// 	cli := redis.NewClusterClient(&redis.ClusterOptions{
// 		Addrs: []string{"localhost:6379"},
// 	})
// 	NewDistLimiterWithRedisCluster(cli)
// 	t1 := time.NewTicker(time.Second)
// 	t2 := time.NewTicker(2 * time.Second)
// 	lim := redis_rate.Limit{
// 		Period: time.Second * 15,
// 		Rate:   5,
// 		Burst:  5,
// 	}
// 	curr := time.Now()
// 	stopat := 2 * time.Minute
// 	for {
// 		select {
// 		case x := <-t1.C:
// 			fmt.Println(int(x.Sub(curr)/time.Second), Limiter(context.Background(), "test", lim))
// 		case y := <-t2.C:
// 			fmt.Println("About to exit")
// 			if y.Sub(curr) >= 2*time.Minute {
// 				t1.Stop()
// 				t2.Stop()
// 			}
// 			break
// 		}
// 		if t1.C == nil && t2.C == nil {
// 			break
// 		}
// 	}
// }
