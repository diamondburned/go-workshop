package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("The time is now", time.Now().Format(time.Kitchen))
	fmt.Println("This slideshow is made in Go!")
}
