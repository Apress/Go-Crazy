package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-queue/nsq"
	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
)

type jobData struct {
	Name    string
	Message string
}

func (j *jobData) Bytes() []byte {
	fmt.Printf("%s:%s\n", j.Name, j.Message)
	res := sleepSomeTime()
	j = &jobData{Name: "I am awake", Message: res}
	b, _ := json.Marshal(j)
	return b
}

func sleepSomeTime() string {
	seconds := rand.Intn(20)
	sleepTime := time.Duration(seconds) * time.Second
	time.Sleep(sleepTime)
	return fmt.Sprintf("Commander, I slept: %d seconds", seconds)
}

func main() {
	rand.Seed(time.Now().Unix())
	taskN := 100
	rets := make(chan string, taskN)

	w := nsq.NewWorker(
		nsq.WithAddr("127.0.0.1:4150"),
		nsq.WithTopic("crazy"),
		nsq.WithChannel("go"),
		nsq.WithMaxInFlight(10),
		nsq.WithRunFunc(func(ctx context.Context, m core.QueuedMessage) error {
			var v *jobData
			if err := json.Unmarshal(m.Bytes(), &v); err != nil {
				return err
			}
			rets <- v.Message
			return nil
		}),
	)

	q := queue.NewPool(10, queue.WithWorker(w))
	defer q.Release()

	for i := 0; i < taskN; i++ {
		go func(i int) {
			q.Queue(&jobData{
				Name:    "Sleeping Gophers",
				Message: fmt.Sprintf("Hello commander, I am handling the job: %d", +i),
			})
		}(i)
	}

	for i := 0; i < taskN; i++ {
		fmt.Println("message:", <-rets)
		time.Sleep(10 * time.Millisecond)
	}
}
