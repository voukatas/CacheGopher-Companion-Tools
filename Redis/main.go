package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	kb        = 1024 // One KB
	totalSize = 64 * kb
	keySize   = 10 // size of the key in bytes
	valueSize = totalSize - keySize - 3
)

// _, err = server.redisClient.Set(ctx, key, questionDataJSON, time.Duration(server.config.QuickGameQuestionStoreTime)*time.Second).Result()
// _, err = server.redisClient.Get(ctx, notifyPoiKey).Result()func main() {
func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // default Redis address
	})

	var durationsSet []time.Duration
	var durationsGet []time.Duration
	value := bytes.Repeat([]byte("x"), valueSize)

	iterations := 10000

	for i := 0; i < iterations; i++ {
		startSet := time.Now()
		_ = rdb.Set(ctx, "testKey", value, 0).Err()
		// if err != nil {
		// 	fmt.Println("Set operation failed:", err)
		// 	continue
		// }
		endSet := time.Since(startSet)
		durationsSet = append(durationsSet, endSet)

		startGet := time.Now()
		_, _ = rdb.Get(ctx, "testKey").Result()
		// if err != nil {
		// 	fmt.Println("Get operation failed:", err)
		// 	continue
		// }
		endGet := time.Since(startGet)
		durationsGet = append(durationsGet, endGet)
		//fmt.Printf("Test value: %s\n", val)
	}

	meanSet := meanDuration(durationsSet)
	meanGet := meanDuration(durationsGet)

	fmt.Printf("Average Set operation took %s\n", meanSet)
	fmt.Printf("Average Get operation took %s\n", meanGet)
}

func meanDuration(durations []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}
