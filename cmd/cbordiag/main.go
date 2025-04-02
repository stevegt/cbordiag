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
	dm, err := cbor.DiagOptions{
		ByteStringText:         true,
		ByteStringEmbeddedCBOR: true,
	}.DiagMode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Diagnose mode: %v\n", err)
		os.Exit(1)
	}
	diagnosis, err := dm.Diagnose(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error diagnosing CBOR data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(diagnosis)

}
