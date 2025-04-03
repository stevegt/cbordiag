package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	data, _ := io.ReadAll(reader)

	if len(data) == 0 {
		fmt.Println("No input data")
		os.Exit(1)
	}

	parser := &cborParser{
		data:   data,
		offset: 0,
		depth:  -1, // Start at -1 to compensate for root level
	}

	var output strings.Builder
	for parser.offset < len(parser.data) {
		lines := parser.parseItem()
		for _, line := range lines {
			output.WriteString(line + "\n")
		}
	}

	fmt.Println(output.String())
}
