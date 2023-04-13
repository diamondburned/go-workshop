package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

// SayHello is a method on Person.
func (p Person) SayHello() {
	fmt.Println("Hello, my name is", p.Name)
}

func sayHelloWithEmojis(h HelloSayer) {
	fmt.Println("👋👋👋")
	h.SayHello()
	fmt.Println("👋👋👋")
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

func main() {
	c := Cat{Name: "Mittens"}
	sayHelloWithEmojis(c)
}
