package main

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func makeRequest(questions []string) gpt3.CompletionRequest {

	maxToken, _ := strconv.Atoi(os.Getenv("MAX_TOKEN"))
	temperature, _ := strconv.ParseFloat(os.Getenv("TEMPERATURE"), 32)

	return gpt3.CompletionRequest{
		Prompt:      questions,
		MaxTokens:   gpt3.IntPtr(maxToken),
		Temperature: gpt3.Float32Ptr(float32(temperature)),
	}
}

func main() {
	programName, questions := os.Args[0], os.Args[1:]
	log.Printf("Starting:%s", programName)

	_ = godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		// exit with message if no key found
		log.Fatalln("Missing ChatGPT API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	request := makeRequest(questions)

	resp, _ := client.CompletionWithEngine(ctx, "davinci", request)

	fmt.Print(resp.Choices[0].Text)
}
