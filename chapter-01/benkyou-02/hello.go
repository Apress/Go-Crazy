package main

import (
	"fmt"
)

type Message struct {
	Hello string
}

func main() {
	h := Message{Hello: "world"}
	fmt.Printf("%+v\n", h)
}
