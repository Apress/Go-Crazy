package main

import (
	"context"
	"fmt"
	"time"
)

func Task(ctx context.Context) {
	var i = 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done")
			return
		default:
			i++
			fmt.Printf("Running [%s]...%d\n", ctx.Value("hello"), i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go Task(context.WithValue(ctx, "hello", "world"))
	go Task(context.WithValue(ctx, "hello", "nico"))

	<-ctx.Done()
}
