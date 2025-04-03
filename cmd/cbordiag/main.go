package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

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

	// Parse entire input and join output lines
	output := strings.Join(cbordiag.Annotate(data), "\n")
	fmt.Println(output)
}
