package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Create("name.txt")
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return
	}
	defer file.Close()

	for i := 0; i < 5; i++ {
		fmt.Fprintf(file, "Hello")
	}
}
