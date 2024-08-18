package main

import (
	"fmt"
	"github.com/lazarus/arithmetic-coding-go/ac"
	"log"
	"os"
	"time"
)

func main() {
	fileName := "ex.txt"
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	context := ac.NewContext()
	start := time.Now()
	encoded := context.Encode(data)
	elapsed := time.Since(start)

	reductionPercentage := float64(len(data)-len(encoded)) / float64(len(data)) * 100

	fmt.Printf("Original size: %d bytes, encoded size: %d bytes, a reduction of %.2f%% in %v!\n", len(data), len(encoded), reductionPercentage, elapsed.String())

	context = ac.NewContext()
	start = time.Now()
	decoded := context.Decode(encoded)
	elapsed = time.Since(start)

	fmt.Printf("Encoded size: %d bytes, decoded size: %d bytes in %v\n", len(encoded), len(decoded), elapsed.String())

	if string(data) == string(decoded) {
		fmt.Println("Success: The original and decoded data match!")
	} else {
		fmt.Println("Error: The original and decoded data do not match.")
	}
}
