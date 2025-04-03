package main

import (
	"fmt"

	"github.com/stevegt/cbordiag"
)

func main() {
	data := []byte{0x83, 0x01, 0x02, 0x03} // CBOR array [1,2,3]
	parser := &cbordiag.CborParser{
		Data:   data,
		Offset: 0,
		Depth:  0,
	}

	diagnostics := parser.Parse()
	for _, line := range diagnostics {
		fmt.Println(line)
	}

	// Output:
	// 83                   # ARRAY (3 items)
	//     01                   # POS INT: 1
	//     02                   # POS INT: 2
	//     03                   # POS INT: 3
}
