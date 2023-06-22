package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch <- fmt.Sprintf("hello")
	}()

	go func() {
		time.Sleep(4 * time.Second)
		ch <- fmt.Sprintf("world")
	}()

	for {
		select {
		case v := <-ch:
			fmt.Printf("%s", v)
		case <-time.After(3 * time.Second):
			fmt.Println("waited 3 seconds")
			os.Exit(0)
		}
	}

}
