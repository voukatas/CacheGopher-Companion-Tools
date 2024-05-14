package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/voukatas/CacheGopher/pkg/client"
)

// a delay is needed because the approach is eventual consistency and for the tests I expect to work as strong consistency so I put the delay
// For 10000 concurrent goroutines that each of them executes a get -> set -> get -> delete -> get scenario, around ~2.5 seconds are required to reach consistency between the servers. This is local test
func main() {
	//conPool := client.NewConnPool(5, "localhost:31337")
	client, err := client.NewClient(false)
	if err != nil {
		fmt.Println("Failed to establish connection:", err.Error())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	numGoroutines := 1

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("testkey%d", id)
			value := fmt.Sprintf("testvalue%d", id)
			//
			if resp, err := client.Get(key); err == nil {
				log.Printf("Get failed!!!: key=%s, expected=%s, got=%s, err=%v", key, value, resp, err)
			}
			//

			if resp, err := client.Set(key, value); err != nil || resp != "OK" {
				log.Printf("Set failed: key=%s, resp=%s, err=%v", key, resp, err)
			}
			//
			// //time.Sleep(2 * time.Second)
			// //time.Sleep(50 * time.Millisecond) //ms
			time.Sleep(2500 * time.Millisecond) //ms
			//
			if resp, err := client.Get(key); err != nil || resp != value {
				log.Printf("Get failed: key=%s, expected=%s, got=%s, err=%v", key, value, resp, err)
			}
			// // time.Sleep(2 * time.Second)
			// //time.Sleep(50 * time.Millisecond) //ms
			time.Sleep(2500 * time.Millisecond) //ms
			//
			if resp, err := client.Delete(key); err != nil || resp != "OK" {
				log.Printf("Delete failed: key=%s, resp=%s, err=%v", key, resp, err)
			}
			// //time.Sleep(2 * time.Second)
			// //time.Sleep(50 * time.Millisecond) //ms
			time.Sleep(2500 * time.Millisecond) //ms
			//
			if resp, err := client.Get(key); err == nil || !strings.Contains(err.Error(), "Key not found") {
				log.Printf("Post-Delete Get should fail: key=%s, resp=%s, err=%v", key, resp, err)
			}
		}(i)
	}

	wg.Wait()

}
