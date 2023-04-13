package main

import (
	"fmt"
	"strings"
)

type Person struct {
	Name string
	Age  int
}

// SayHello is a method on Person.
func (p Person) SayHello() {
	fmt.Println("Hello, my name is", p.Name)
}

type Cat struct {
	Name string
}

func (c Cat) SayHello() {
	fmt.Println("Meow, meow name is", c.Name)
}

type HelloSayer interface {
	SayHello()
}

func sayHelloWithEmojis(h HelloSayer) {
	emoji := "👋"
	// use a cat emoji if h is a cat
	if _, ok := h.(Cat); ok {
		emoji = "🐾"
	}

	fmt.Println(strings.Repeat(emoji, 3))
	h.SayHello()
	fmt.Println(strings.Repeat(emoji, 3))
}

func main() {
	sayHelloWithEmojis(Cat{Name: "Mittens"})
	sayHelloWithEmojis(Person{Name: "John", Age: 20})
}
