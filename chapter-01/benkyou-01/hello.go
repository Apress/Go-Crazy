package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, _ := os.OpenFile("hello.txt", os.O_RDONLY, 0666)
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		fmt.Printf("> %s", line)
		if err != nil {
			return
		}
	}
}
