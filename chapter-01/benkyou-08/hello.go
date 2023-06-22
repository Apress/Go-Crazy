package main

import (
	"fmt"
	"time"
)

func printNumbers(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i
		time.Sleep(100 * time.Millisecond)
	}
	close(c)
}

func main() {
	c := make(chan int)
	go printNumbers(c)

	for num := range c {
		fmt.Println(num)
	}
}
