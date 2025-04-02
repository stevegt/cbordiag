package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type cborParser struct {
	data   []byte
	offset int
	depth  int
}

func (p *cborParser) parseItem() string {
	if p.offset >= len(p.data) {
		return ""
	}

	initial := p.data[p.offset]
	p.offset++
	major := initial >> 5
	info := initial & 0x1F

	var annotation string
	var value string

	switch major {
	case 0: // Unsigned integer
		val := p.parseUint(info)
		annotation = fmt.Sprintf("POS INT: %d", val)
	case 2: // Byte string
		length := p.parseLength(info)
		bytes := p.data[p.offset : p.offset+length]
		if isPrintable(bytes) {
			value = fmt.Sprintf("%q", bytes)
		} else {
			value = fmt.Sprintf("h'%x'", bytes)
		}
		p.offset += length
		annotation = fmt.Sprintf("BYTE STR: %s (%d bytes)", value, length)
	case 3: // Text string
		length := p.parseLength(info)
		str := string(p.data[p.offset : p.offset+length])
		p.offset += length
		annotation = fmt.Sprintf("TEXT: %q (%d bytes)", str, length)
	case 4: // Array
		length := p.parseLength(info)
		annotation = fmt.Sprintf("ARRAY (%d items)", length)
		p.depth++
		items := p.parseCollection(length)
		p.depth--
		value = "\n" + items
	case 5: // Map
		pairs := p.parseLength(info)
		annotation = fmt.Sprintf("MAP (%d pairs)", pairs)
		p.depth++
		items := p.parseCollection(pairs*2)
		p.depth--
		value = "\n" + items
	default:
		annotation = "UNKNOWN TYPE"
	}

	prefix := fmt.Sprintf("%02X", initial)
	if info > 0x17 {
		additional := p.data[p.offset-1-p.bytesUsed(info): p.offset]
		prefix = fmt.Sprintf("%02X", initial) + " " + strings.ToUpper(hex.EncodeToString(additional))
	}

	return fmt.Sprintf("%-20s # %s%s", prefix, annotation, value)
}

func (p *cborParser) parseUint(info byte) uint64 {
	// Implementation for parsing unsigned integers
}

func (p *cborParser) parseLength(info byte) int {
	// Implementation for length parsing
}

func (p *cborParser) parseCollection(count int) string {
	// Implementation for nested structures
}

func (p *cborParser) bytesUsed(info byte) int {
	// Implementation for byte consumption calculation
}

func isPrintable(b []byte) bool {
	// Implementation for printable check
}

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
		depth:  0,
	}

	var output strings.Builder
	for parser.offset < len(parser.data) {
		line := parser.parseItem()
		output.WriteString(strings.Repeat("|-- ", parser.depth) + line + "\n")
	}

	fmt.Println(output.String())
}
