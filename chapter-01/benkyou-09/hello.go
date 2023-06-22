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

	for value := range c {
		switch value % 2 {
		case 0:
			fmt.Printf("Value: %d is even\n", value)
		case 1:
			fmt.Printf("Value: %d is odd\n", value)
		default:
			fmt.Println("Received a weird value")
		}
	}
}
