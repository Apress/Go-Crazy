package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	for true {
		fmt.Print("\n\n> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("read line: %s-\n", line)
		complete(ctx, client, line)
	}
}

func makeRequest(question string) gpt3.CompletionRequest {

	maxToken, _ := strconv.Atoi(os.Getenv("MAX_TOKEN"))
	temperature, _ := strconv.ParseFloat(os.Getenv("TEMPERATURE"), 32)

	questions := []string{question}
	return gpt3.CompletionRequest{
		Prompt:      questions,
		MaxTokens:   gpt3.IntPtr(maxToken),
		Temperature: gpt3.Float32Ptr(float32(temperature)),
	}
}

func complete(ctx context.Context, client gpt3.Client, question string) {

	request := makeRequest(question)
	resp, _ := client.Completion(ctx, request)

	fmt.Print(resp.Choices[0].Text)
}
