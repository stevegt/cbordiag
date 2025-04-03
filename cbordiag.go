package cbordiag

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
)

type CborParser struct {
	Data   []byte
	Offset int
	Depth  int
}

func (p *CborParser) ParseItem() []string {
	if p.Offset >= len(p.Data) {
		return []string{}
	}

	startOffset := p.Offset
	initial := p.Data[p.Offset]
	p.Offset++
	major := initial >> 5
	info := initial & 0x1F

	var lines []string
	var annotation, value string
	var length int

	// Helper to create error lines with hex prefix
	errorLine := func(msg string, args ...interface{}) []string {
		prefixBytes := p.Data[startOffset:p.Offset]
		prefix := strings.ToUpper(hex.EncodeToString(prefixBytes))
		indent := strings.Repeat("    ", p.Depth)
		errMsg := fmt.Sprintf(msg, args...)
		return []string{fmt.Sprintf("%s%-20s # ERROR: %s", indent, prefix, errMsg)}
	}

	switch major {
	case 0: // Unsigned integer
		val := p.parseUint(info)
		annotation = fmt.Sprintf("POS INT: %d", val)
	case 1: // Negative integer
		val := p.parseNint(info)
		annotation = fmt.Sprintf("NEG INT: %d", val)
	case 2: // Byte string
		length = p.parseLength(info)
		if p.Offset+length > len(p.Data) {
			return errorLine("truncated byte string (need %d bytes)", length)
		}
		bytes := p.Data[p.Offset : p.Offset+length]
		p.Offset += length
		if isPrintable(bytes) {
			value = fmt.Sprintf("'%s'", bytes)
		} else {
			value = fmt.Sprintf("h'%x'", bytes)
		}
		annotation = fmt.Sprintf("BYTE STR: %s (%d bytes)", value, length)
	case 3: // Text string
		length = p.parseLength(info)
		if p.Offset+length > len(p.Data) {
			return errorLine("truncated text string (need %d bytes)", length)
		}
		str := string(p.Data[p.Offset : p.Offset+length])
		p.Offset += length
		annotation = fmt.Sprintf("TEXT: %q (%d bytes)", str, length)
	case 4: // Array
		length = p.parseLength(info)
		annotation = fmt.Sprintf("ARRAY (%d items)", length)
	case 5: // Map
		length = p.parseLength(info)
		annotation = fmt.Sprintf("MAP (%d pairs)", length)
	case 6: // Tag
		tag := p.parseUint(info)
		annotation = fmt.Sprintf("TAG (%d)", tag)
	case 7: // Special
		if info < 20 {
			annotation = fmt.Sprintf("SIMPLE: %d", info)
		} else {
			annotation = "SIMPLE (RESERVED)"
		}
	}

	// Handle nesting for composite types
	if major == 4 || major == 5 || major == 6 {
		p.Depth++
		defer func() { p.Depth-- }()
	}

	// Create main annotation line
	prefixBytes := p.Data[startOffset:p.Offset]
	prefix := strings.ToUpper(hex.EncodeToString(prefixBytes))
	indent := strings.Repeat("    ", p.Depth)
	lines = append(lines, fmt.Sprintf("%s%-20s # %s", indent, prefix, annotation))

	// Process nested items
	switch major {
	case 4: // Array
		for i := 0; i < length; i++ {
			if p.Offset >= len(p.Data) {
				lines = append(lines, errorLine("truncated array")...)
				break
			}
			lines = append(lines, p.ParseItem()...)
		}
	case 5: // Map
		for i := 0; i < length*2; i++ {
			if p.Offset >= len(p.Data) {
				lines = append(lines, errorLine("truncated map")...)
				break
			}
			lines = append(lines, p.ParseItem()...)
		}
	case 6: // Tag
		if p.Offset < len(p.Data) {
			lines = append(lines, p.ParseItem()...)
		}
	}

	return lines
}

func (p *CborParser) parseUint(info byte) uint64 {
	if info < 24 {
		return uint64(info)
	}
	var val uint64
	switch info {
	case 24:
		if p.Offset+1 > len(p.Data) {
			return 0
		}
		val = uint64(p.Data[p.Offset])
		p.Offset++
	case 25:
		if p.Offset+2 > len(p.Data) {
			return 0
		}
		val = uint64(p.Data[p.Offset])<<8 | uint64(p.Data[p.Offset+1])
		p.Offset += 2
	case 26:
		if p.Offset+4 > len(p.Data) {
			return 0
		}
		val = uint64(p.Data[p.Offset])<<24 | uint64(p.Data[p.Offset+1])<<16 |
			uint64(p.Data[p.Offset+2])<<8 | uint64(p.Data[p.Offset+3])
		p.Offset += 4
	case 27:
		if p.Offset+8 > len(p.Data) {
			return 0
		}
		val = uint64(p.Data[p.Offset])<<56 | uint64(p.Data[p.Offset+1])<<48 |
			uint64(p.Data[p.Offset+2])<<40 | uint64(p.Data[p.Offset+3])<<32 |
			uint64(p.Data[p.Offset+4])<<24 | uint64(p.Data[p.Offset+5])<<16 |
			uint64(p.Data[p.Offset+6])<<8 | uint64(p.Data[p.Offset+7])
		p.Offset += 8
	}
	return val
}

func (p *CborParser) parseNint(info byte) int64 {
	return -1 - int64(p.parseUint(info))
}

func (p *CborParser) parseLength(info byte) int {
	if info < 24 {
		return int(info)
	}
	var length int
	switch info {
	case 24:
		if p.Offset+1 > len(p.Data) {
			return 0
		}
		length = int(p.Data[p.Offset])
		p.Offset++
	case 25:
		if p.Offset+2 > len(p.Data) {
			return 0
		}
		length = int(p.Data[p.Offset])<<8 | int(p.Data[p.Offset+1])
		p.Offset += 2
	case 26:
		if p.Offset+4 > len(p.Data) {
			return 0
		}
		length = int(p.Data[p.Offset])<<24 | int(p.Data[p.Offset+1])<<16 |
			int(p.Data[p.Offset+2])<<8 | int(p.Data[p.Offset+3])
		p.Offset += 4
	case 27:
		if p.Offset+8 > len(p.Data) {
			return 0
		}
		length = int(uint64(p.Data[p.Offset])<<56 | uint64(p.Data[p.Offset+1])<<48 |
			uint64(p.Data[p.Offset+2])<<40 | uint64(p.Data[p.Offset+3])<<32 |
			uint64(p.Data[p.Offset+4])<<24 | uint64(p.Data[p.Offset+5])<<16 |
			uint64(p.Data[p.Offset+6])<<8 | uint64(p.Data[p.Offset+7]))
		p.Offset += 8
	}
	return length
}

func isPrintable(b []byte) bool {
	for _, c := range b {
		if c > unicode.MaxASCII || !unicode.IsPrint(rune(c)) {
			return false
		}
	}
	return true
}
