package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/voukatas/CacheGopher/pkg/client"
)

const (
	kb        = 1024
	totalSize = 64 * kb
	keySize   = 10
	valueSize = totalSize - keySize - 3
)

func main() {
	//conPool := client.NewConnPool(5, "localhost:31337")
	client, err := client.NewClient(false)
	if err != nil {
		fmt.Println("Failed to establish connection:", err.Error())
		os.Exit(1)
	}

	value := strings.Repeat("x", valueSize)
	//fmt.Println(value)

	var durationsSet []time.Duration
	var durationsGet []time.Duration
	iterations := 1000

	for i := 0; i < iterations; i++ {
		startSet := time.Now()
		_, err = client.Set("testKey", value)
		// if err != nil {
		// 	fmt.Println("Set operation failed:", err)
		// 	continue
		// }
		endSet := time.Since(startSet)
		durationsSet = append(durationsSet, endSet)

		startGet := time.Now()
		client.Get("testKey")
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
