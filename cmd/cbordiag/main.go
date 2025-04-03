package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/stevegt/cbordiag"
)

// main reads CBOR data from stdin and outputs annotated diagnostic format
func main() {
	reader := bufio.NewReader(os.Stdin)
	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if len(data) == 0 {
		fmt.Println("No input data")
		os.Exit(1)
	}

	parser := &cbordiag.CborParser{
		Data:   data,
		Offset: 0,
		Depth:  0,
	}

	// Parse entire input and join output lines
	output := strings.Join(parser.Parse(), "\n")
	fmt.Println(output)
}
