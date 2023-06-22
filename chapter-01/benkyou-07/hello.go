package main

import (
	"fmt"
	"os"
)

func main() {
	programName, questions := os.Args[0], os.Args[1:]
	fmt.Printf("Starting:%s", programName)
	if len(questions) == 0 {
		fmt.Printf("Usage:%s <question1> <question2> …", programName)
	} else {
		for i, question := range questions {
			fmt.Printf("Question [%d] > %s\n", i, question)
		}
	}

}
