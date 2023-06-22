package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Task finished")
		cancel()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Context Done")
		err := ctx.Err()
		if err != nil {
			fmt.Printf("err: %s", err)
		}

	}
}
