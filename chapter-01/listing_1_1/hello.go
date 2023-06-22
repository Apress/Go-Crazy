package main

import (
	"fmt"
	"runtime"
)

func main() {
	version := runtime.Version()
	fmt.Printf("Go version: %s\n", version)
}
