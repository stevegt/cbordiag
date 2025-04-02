# CBOR Diagnostic Tool

A command-line tool for CBOR diagnostics with annotated output, inspired by RFC 8949 EDN and RFC 8610 extensions. Implements a zero-dependency CBOR parser with type annotations.

## Features

- Input formats:
  - Raw binary CBOR
  - Hexadecimal strings
  - Standard input streams
- Diagnostic features:
  - Annotated hex output with type metadata
  - Major type breakdown (RFC 8949 §3)
  - Embedded CBOR decoding (RFC 8610 §G.3)
  - Nested structure visualization

## Installation

```bash
git clone https://github.com/yourusername/cbordiag.git
cd cbordiag
go build -o cbordiag ./cmd/cbordiag
```

## Usage

```bash
# Process hex input
echo a26161016162820203 | ./cbordiag

# Read binary file 
./cbordiag < data.cbor

# Validate dCBOR
cbor-generator | ./cbordiag -v
```

## Annotation Example

Input (hex):
```
a26161016162820203
```

Output:
```text
a2                 # MAP (2 pairs)
   61              # TEXT (1 byte)
      61           # "a"
   01              # UINT (1)
   62              # TEXT (1 byte)
      62           # "b"
   82              # ARRAY (2 items)
      02           # UINT (2)
      03           # UINT (3)
```

## Extended Features

```bash
# Show hex annotations
echo a16568656c6c6f65776f726c64 | ./cbordiag -a hex
```

Output:
```text
a1                 # MAP (1 pair)
   65 68656c6c6f   # TEXT: "hello" (5 bytes)
   65 776f726c64   # TEXT: "world" (5 bytes)
```

## Options

| Flag | Description                |
|------|----------------------------|
| -i   | Input format (bin/hex/auto)|
| -a   | Annotation depth (1-99)    |
