# CBOR Diagnostic Tool

A command-line tool for diagnosing CBOR data with annotated output. Supports binary, hex, and standard input formats with extended diagnostic explanations.

## Features

- Input formats supported:
  - Raw binary CBOR
  - Hexadecimal encoded CBOR
  - Standard input streams
- Output features:
  - Extended Diagnostic Notation (EDN) with type annotations
  - Structured breakdown of CBOR major types
  - Embedded CBOR decoding for nested structures
- Validation of Deterministic CBOR (dCBOR)

## Installation

```bash
git clone https://github.com/yourusername/cbordiag.git
cd cbordiag
go build -o cbordiag ./cmd/cbordiag
```

## Usage 

```bash
# Read from stdin with hex input
echo a16568656c6c6f65776f726c64 | ./cbordiag

# Read binary file
./cbordiag < data.cbor

# Piped conversion example
cbor-generator | ./cbordiag
```

## Annotated Output Example

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

| Flag | Description               |
|------|---------------------------|
| -a   | Annotation style (compact/verbose) |
| -h   | Show help                 |
