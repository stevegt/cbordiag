package cbordiag

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
)

// CborParser maintains state for parsing CBOR data and generating diagnostics
type CborParser struct {
	Data   []byte // Input CBOR data
	Offset int    // Current read position in Data
	Depth  int    // Current nesting depth for indentation
}

// Parse processes all CBOR items starting at current offset and returns
// formatted diagnostic lines with proper indentation and annotations
func (p *CborParser) Parse() []string {
	var lines []string
	for p.Offset < len(p.Data) {
		lines = append(lines, p.ParseItem()...)
	}
	return lines
}

// ParseItem parses a single CBOR item starting at current offset and returns
// the formatted diagnostic lines with proper indentation. Handles nested
// structures by recursively parsing contained items.
func (p *CborParser) ParseItem() []string {
	if p.Offset >= len(p.Data) {
		return []string{}
	}

	startOffset := p.Offset
	initial := p.Data[p.Offset]
	p.Offset++
	major := initial >> 5  // Extract major type (high 3 bits)
	info := initial & 0x1F // Extract additional info (low 5 bits)

	var lines []string
	var annotation, value string
	var length int

	// errorLine generates an error message with truncated hex prefix
	errorLine := func(msg string, args ...interface{}) []string {
		prefixBytes := p.Data[startOffset:]
		if len(prefixBytes) > 4 { // Truncate long prefixes for errors
			prefixBytes = prefixBytes[:4]
		}
		prefix := strings.ToUpper(hex.EncodeToString(prefixBytes))
		indent := strings.Repeat("    ", p.Depth)
		errMsg := fmt.Sprintf(msg, args...)
		return []string{fmt.Sprintf("%s%-20s # ERROR: %s", indent, prefix, errMsg)}
	}

	// Handle major types according to CBOR specification
	switch major {
	case 0: // Unsigned integer (RFC 8949 Section 3.1)
		val := p.parseUint(info)
		annotation = fmt.Sprintf("POS INT: %d", val)
	case 1: // Negative integer (RFC 8949 Section 3.1)
		val := p.parseNint(info)
		annotation = fmt.Sprintf("NEG INT: %d", val)
	case 2: // Byte string (RFC 8949 Section 3.2.1)
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
		annotation = fmt.Sprintf("BYTE STR: %s (%d %s)", 
			value, length, pluralize(length, "byte", "bytes"))
	case 3: // Text string (RFC 8949 Section 3.2.1)
		length = p.parseLength(info)
		if p.Offset+length > len(p.Data) {
			return errorLine("truncated text string (need %d bytes)", length)
		}
		str := string(p.Data[p.Offset : p.Offset+length])
		p.Offset += length
		annotation = fmt.Sprintf("TEXT: %q (%d %s)", 
			str, length, pluralize(length, "byte", "bytes"))
	case 4: // Array (RFC 8949 Section 3.2.2)
		length = p.parseLength(info)
		annotation = fmt.Sprintf("ARRAY (%d %s)", 
			length, pluralize(length, "item", "items"))
	case 5: // Map (RFC 8949 Section 3.2.3)
		length = p.parseLength(info)
		annotation = fmt.Sprintf("MAP (%d %s)", 
			length, pluralize(length, "pair", "pairs"))
	case 6: // Tag (RFC 8949 Section 3.4)
		tag := p.parseUint(info)
		annotation = fmt.Sprintf("TAG (%d)", tag)
	case 7: // Special values (RFC 8949 Section 3.3)
		if info < 20 {
			annotation = fmt.Sprintf("SIMPLE: %d", info)
		} else {
			annotation = "SIMPLE (RESERVED)"
		}
	}

	// Generate prefix from consumed bytes
	prefixBytes := p.Data[startOffset:p.Offset]
	prefix := strings.ToUpper(hex.EncodeToString(prefixBytes))
	indent := strings.Repeat("    ", p.Depth)
	lines = append(lines, fmt.Sprintf("%s%-20s # %s", indent, prefix, annotation))

	// Process nested structures recursively
	if major == 4 || major == 5 || major == 6 {
		p.Depth++
		defer func() { p.Depth-- }()

		switch major {
		case 4: // Array elements
			for i := 0; i < length; i++ {
				if p.Offset >= len(p.Data) {
					lines = append(lines, errorLine("truncated array")...)
					break
				}
				lines = append(lines, p.ParseItem()...)
			}
		case 5: // Map key-value pairs
			for i := 0; i < length*2; i++ {
				if p.Offset >= len(p.Data) {
					lines = append(lines, errorLine("truncated map")...)
					break
				}
				lines = append(lines, p.ParseItem()...)
			}
		case 6: // Tagged item
			if p.Offset < len(p.Data) {
				lines = append(lines, p.ParseItem()...)
			}
		}
	}

	return lines
}

// parseUint decodes CBOR unsigned integers with error checking
func (p *CborParser) parseUint(info byte) uint64 {
	if info < 24 { // Direct value (RFC 8949 Section 3.1)
		return uint64(info)
	}
	
	var val uint64
	switch info {
	case 24: // 1-byte value
		if p.Offset+1 > len(p.Data) {
			return 0
		}
		val = uint64(p.Data[p.Offset])
		p.Offset++
	case 25: // 2-byte value
		if p.Offset+2 > len(p.Data) {
			return 0
		}
		val = uint64(p.Data[p.Offset])<<8 | uint64(p.Data[p.Offset+1])
		p.Offset += 2
	case 26: // 4-byte value
		if p.Offset+4 > len(p.Data) {
			return 0
		}
		val = uint64(p.Data[p.Offset])<<24 | uint64(p.Data[p.Offset+1])<<16 |
			uint64(p.Data[p.Offset+2])<<8 | uint64(p.Data[p.Offset+3])
		p.Offset += 4
	case 27: // 8-byte value
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

// parseNint calculates negative integers using n = -1 - val encoding
func (p *CborParser) parseNint(info byte) int64 {
	return -1 - int64(p.parseUint(info))
}

// parseLength decodes length values for CBOR data structures
func (p *CborParser) parseLength(info byte) int {
	if info < 24 { // Direct length
		return int(info)
	}

	var length int
	switch info {
	case 24: // 1-byte length
		if p.Offset+1 > len(p.Data) {
			return 0
		}
		length = int(p.Data[p.Offset])
		p.Offset++
	case 25: // 2-byte length
		if p.Offset+2 > len(p.Data) {
			return 0
		}
		length = int(p.Data[p.Offset])<<8 | int(p.Data[p.Offset+1])
		p.Offset += 2
	case 26: // 4-byte length
		if p.Offset+4 > len(p.Data) {
			return 0
		}
		length = int(p.Data[p.Offset])<<24 | int(p.Data[p.Offset+1])<<16 |
			int(p.Data[p.Offset+2])<<8 | int(p.Data[p.Offset+3])
		p.Offset += 4
	case 27: // 8-byte length
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

// isPrintable checks if bytes contain only printable ASCII characters
func isPrintable(b []byte) bool {
	for _, c := range b {
		if c > unicode.MaxASCII || !unicode.IsPrint(rune(c)) {
			return false
		}
	}
	return true
}

// pluralize returns singular or plural form based on count
func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
