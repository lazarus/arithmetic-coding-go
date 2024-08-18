package main

import (
	"arithmetic-coding/ac"
	"fmt"
	"log"
	"os"
)

func main() {
	fileName := "ex.txt"
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	context := ac.NewContext()
	encoded := context.Encode(data)

	fmt.Printf("Original size: %d bytes, encoded size: %d bytes\n", len(data), len(encoded))

	context = ac.NewContext()
	decoded := context.Decode(encoded)

	fmt.Printf("Encoded size: %d bytes, decoded size: %d bytes\n", len(encoded), len(decoded))

	if string(data) == string(decoded) {
		fmt.Println("Success: The original and decoded data match!")
	} else {
		fmt.Println("Error: The original and decoded data do not match.")
	}
}
