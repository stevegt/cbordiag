package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fxamacker/cbor/v2"
)

func main() {
	// Create a buffered reader for stdin
	reader := bufio.NewReader(os.Stdin)

	// Read all data from stdin
	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// Diagnose the CBOR data
	diagnosis, err := cbor.Diagnose(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error diagnosing CBOR data: %v\n", err)
		os.Exit(1)
	}

	// Print the diagnosis to stdout
	fmt.Println(diagnosis)
}
