# Adaptive Arithmetic Coding in Go

This repository contains a Go implementation of adaptive arithmetic coding, specifically utilizing a variant of the method described by Tjalling J. Willems.

## Features

- **Adaptive Model**: Frequencies are updated on-the-fly, allowing the model to adapt to the data being processed.
- **Binary Arithmetic Coding**: Operates at the bit level, providing high compression efficiency.
- **Willems' Variant**: Uses a frozen state to prevent frequency overflow, ensuring stable and accurate compression over long data streams.

## Installation

To install this package, you can simply clone the repository:

```bash
git clone https://github.com/lazarus/arithmetic-coding-go.git
```

Or use Go modules to include it in your project:

```bash
go get github.com/lazarus/arithmetic-coding-go
```

## Usage

### Basic Example

Here’s a simple example of how to use the library to encode and decode a file:

```go
package main

import (
"fmt"
"log"
"os"
"arithmetic-coding-go/ac" // Replace with your actual import path
)

func main() {
// Read input file
fileName := "ex.txt"
data, err := os.ReadFile(fileName)
if err != nil {
log.Fatalf("Failed to read file: %v", err)
}

    // Create a new arithmetic coding context
    context := ac.NewContext()
    
    // Encode the file data
    encoded := context.Encode(data)
    fmt.Printf("Original size: %d bytes, Encoded size: %d bytes\n", len(data), len(encoded))

    // Decode the encoded data
    context = ac.NewContext()
    decoded := context.Decode(encoded)
    fmt.Printf("Encoded size: %d bytes, Decoded size: %d bytes\n", len(encoded), len(decoded))

    // Check if the original and decoded data match
    if string(data) == string(decoded) {
        fmt.Println("Success: The original and decoded data match!")
    } else {
        fmt.Println("Error: The original and decoded data do not match.")
    }
}
```

### API Overview

- **NewContext**: Initializes a new encoding/decoding context.
- **Encode**: Encodes a slice of bytes into a compressed format.
- **Decode**: Decodes a previously encoded slice of bytes back into the original data.

### Example File

Place your input file (`ex.txt`) in the same directory as your Go program, and run the above example to see the encoding and decoding process in action.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributions

Contributions are welcome! Feel free to open an issue or submit a pull request.

## Acknowledgments

This implementation is based on the adaptive arithmetic coding method described by Tjalling J. Willems.
