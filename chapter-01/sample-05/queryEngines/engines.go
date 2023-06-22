package main

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	engines, err := client.Engines(ctx)
	if err != nil {
		return
	}

	for _, engine := range engines.Data {
		fmt.Printf("Engine ID: %s, Name: %s, Ready: %t\n", engine.ID, engine.Owner, engine.Ready)
	}

}
