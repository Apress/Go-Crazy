package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-queue/queue"
)

func sleepSomeTime() string {
	sleepTime := time.Duration(rand.Intn(60)) * time.Second
	message := fmt.Sprintf("%s\n", sleepTime)
	fmt.Printf("About to process: %s\n", message)
	time.Sleep(sleepTime)
	return message
}

func job(i int, rets chan string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		sleepSomeTime()
		rets <- fmt.Sprintf("Hello commander, I am handling the job: %02d", +i)
		return nil
	}
}

func main() {
	taskN := 100
	rets := make(chan string, taskN)

	q := queue.NewPool(5)
	defer q.Release()

	for i := 0; i < taskN; i++ {
		go q.QueueTask(job(i, rets))
	}

	for i := 0; i < taskN; i++ {
		fmt.Println("message:", <-rets)
		time.Sleep(20 * time.Millisecond)
	}
}
