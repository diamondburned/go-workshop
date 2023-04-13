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

// BEGIN OMIT
// HelloSayer describes anything that can say hello.
// It's conventional to name interfaces with the -er suffix.
type HelloSayer interface {
	SayHello()
}

func sayHelloWithEmojis(h HelloSayer) {
	fmt.Println("👋👋👋")
	h.SayHello()
	fmt.Println("👋👋👋")
}

func main() {
	p := Person{Name: "John", Age: 20}
	sayHelloWithEmojis(p)
}

// END OMIT
