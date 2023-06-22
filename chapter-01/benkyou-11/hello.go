package main

import (
	"context"
	"fmt"
	"time"
)

func runTask(ctx context.Context) error {
	if data, ok := ctx.Value("data").(string); ok {
		fmt.Println("Data passed to context:", data)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(2 * time.Second):
		fmt.Println("Task finished")
		return nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), "data", "hello world"), 1*time.Second)
	defer cancel()

	if err := runTask(ctx); err != nil {
		fmt.Println("Task cancelled:", err)
	}
}
