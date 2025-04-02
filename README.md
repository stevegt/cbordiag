# CBOR Diagnostic Tool

A simple command-line tool for diagnosing CBOR (Concise Binary Object Representation) data. This tool decodes CBOR input and provides a human-readable representation of the data structure with annotations explaining each element.

## Features

- Reads CBOR data from standard input (stdin).
- Decodes the CBOR data using the `github.com/fxamacker/cbor/v2` library's `Diagnose()` function.
- Outputs the decoded result with annotations to standard output (stdout).

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or later installed on your system.

### Build

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/cbordiag.git
   ```

2. Navigate to the project directory:

   ```bash
   cd cbordiag
   ```

3. Build the executable:

   ```bash
   go build -o cbordiag ./cmd/cbordiag
   ```

   This will generate the `cbordiag` executable in the project root.

## Usage

You can use the `cbordiag` tool by piping CBOR data into it. Here's a simple example:

```bash
cat yourdata.cbor | ./cbordiag
```

Alternatively, you can use it with other commands that generate CBOR data:

```bash
echo "your cbor data" | ./cbordiag
```

## Example

Suppose you have a CBOR-encoded file named `data.cbor`:

```bash
cat data.cbor | ./cbordiag
```

The tool will output a human-readable diagnosis of the CBOR data structure with annotations. For example:

```
0(foo)       # Major Type 0: Unsigned Integer or Simple Value
1(bar)       # Major Type 1: Negative Integer
2(42)        # Major Type 2: Byte String
3(text)      # Major Type 3: UTF-8 String
4(array)     # Major Type 4: Array
5(map)       # Major Type 5: Map
6(tagged)    # Major Type 6: Tag
7(simple)    # Major Type 7: Floating Point Numbers and Simple Data Types
```

Each line represents a CBOR data element with an annotation explaining its type.

## Dependencies

- [fxamacker/cbor](https://github.com/fxamacker/cbor) v2

## License

MIT License. See [LICENSE](LICENSE) for details.
