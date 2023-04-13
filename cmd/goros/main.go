package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	defer func() { fmt.Println("Took", time.Since(start)) }()

	var wg sync.WaitGroup // counter that blocks until it reaches zero

	for i := 0; i < 5; i++ {
		wg.Add(1)        // count our goroutine
		go func(i int) { // start the goroutine ("go"!)
			defer wg.Done() // uncount when we're done

			time.Sleep(1 * time.Second)
			fmt.Printf("Hello from goroutine %d!\n", i)
		}(i)
	}

	wg.Wait()
}
