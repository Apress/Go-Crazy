package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Hello   string `json:"hellooo"`
	ignored string
}

func main() {
	h := Message{Hello: "world"}
	AsString, _ := json.Marshal(h)
	fmt.Printf("%s\n", AsString)
}
