# CBOR Diagnostic Tool

A Go implementation for parsing CBOR (Concise Binary Object Representation) data and generating annotated diagnostic output with nested structure visualization.

## Features

- Detailed annotation of CBOR data types and values
- Recursive parsing of nested structures (arrays/maps/tags)
- Error detection for truncated data
- Printable detection for byte strings
- Standard compliance with RFC 8949 specifications

## Quickstart

### Installation
```bash
go install github.com/stevegt/cbordiag/cmd/cbordiag@latest
```

### CLI Usage
Process CBOR data from stdin:
```bash
# Process hex-encoded CBOR input
echo -n "8301020382040f" | xxd -r -p | cbordiag

# Process raw CBOR file
cbordiag < data.cbor
```

Sample output for `echo -n "8301020382040f" | xxd -r -p | cbordiag`:
```
83                   # ARRAY (3 items)
    01                   # POS INT: 1
    02                   # POS INT: 2
    03                   # POS INT: 3
82040F              # ARRAY (2 items)
    04                   # POS INT: 4
    0F                   # POS INT: 15
```

### Library Usage
```go
package main

import (
    "fmt"
    "github.com/stevegt/cbordiag"
)

func main() {
    data := []byte{0x83, 0x01, 0x02, 0x03} // CBOR array [1,2,3]
    diagnostics := cbordiag.Annotate(data)
    for _, line := range diagnostics {
        fmt.Println(line)
    }
    
    // Output:
    // 83                   # ARRAY (3 items)
    //     01                   # POS INT: 1
    //     02                   # POS INT: 2
    //     03                   # POS INT: 3
}
```

## Diagnostic Format
Each line follows this structure:
```
<HEX BYTES>         # <TYPE ANNOTATION> [(<DETAILS>)]
```
- **Hex Bytes**: Raw CBOR bytes in uppercase hex
- **Type Annotation**: CBOR major type and value details
- **Indentation**: 4 spaces per nesting level
- **Error Handling**: Truncated data marked with `ERROR`

## Testing
Run comprehensive tests:
```bash
go test -v ./...
```

## References
- RFC 8949: Concise Binary Object Representation (CBOR)

