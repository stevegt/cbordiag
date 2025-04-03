package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/stevegt/cbordiag"
)

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

	var output strings.Builder
	for parser.Offset < len(parser.Data) {
		lines := parser.ParseItem()
		for _, line := range lines {
			output.WriteString(line + "\n")
		}
	}

	fmt.Print(output.String())
}
