package main

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
)

func main() {
	apiKey := "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"
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	request := gpt3.CompletionRequest{
		Prompt: []string{"How many cups of coffee should I drink per day?"},
	}
	resp, _ := client.Completion(ctx, request)

	fmt.Print(resp.Choices[0].Text)
}
