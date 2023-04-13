package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex

	mu.Lock()
	defer mu.Unlock()

	fmt.Print("Hello")
	defer fmt.Println("world")
	fmt.Print(", ")

	// Actual order of execution:
	// 0. mu.Lock()
	// 1. fmt.Print("Hello")
	// 2. fmt.Print(", ")
	// 3. fmt.Println("world")
	// 4. mu.Unlock()
}
