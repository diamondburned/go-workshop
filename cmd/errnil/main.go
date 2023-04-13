package main

import (
	"fmt"
	"log"
	"strconv"
)

func main() {
	if err := do(); err != nil {
		log.Fatalln(err, ":(")
	}
}

func do() error {
	v, err := strconv.Atoi("pee poo")
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}
