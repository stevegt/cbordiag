package cbordiag

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
)

type cborParser struct {
	data   []byte
	offset int
	depth  int
}

func (p *cborParser) parseItem() []string {
	if p.offset >= len(p.data) {
		return []string{}
	}

	startOffset := p.offset
	initial := p.data[p.offset]
	p.offset++
	major := initial >> 5
	info := initial & 0x1F

	var lines []string
	var annotation, value string
	var length int

	switch major {
	case 0: // Unsigned integer
		val := p.parseUint(info)
		annotation = fmt.Sprintf("POS INT: %d", val)
	case 1: // Negative integer
		val := p.parseNint(info)
		annotation = fmt.Sprintf("NEG INT: %d", val)
	case 2: // Byte string
		length = p.parseLength(info)
		if p.offset+length > len(p.data) {
			return []string{fmt.Sprintf("%sERROR: truncated byte string (need %d bytes)",
				strings.Repeat("    ", p.depth), length)}
		}
		bytes := p.data[p.offset : p.offset+length]
		p.offset += length
		if isPrintable(bytes) {
			value = fmt.Sprintf("'%s'", bytes)
		} else {
			value = fmt.Sprintf("h'%x'", bytes)
		}
		annotation = fmt.Sprintf("BYTE STR: %s (%d bytes)", value, length)
	case 3: // Text string
		length = p.parseLength(info)
		if p.offset+length > len(p.data) {
			return []string{fmt.Sprintf("%sERROR: truncated text string (need %d bytes)",
				strings.Repeat("    ", p.depth), length)}
		}
		str := string(p.data[p.offset : p.offset+length])
		p.offset += length
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

	if major == 4 || major == 5 || major == 6 {
		p.depth++
		defer func() { p.depth-- }()
	}

	prefixBytes := p.data[startOffset:p.offset]
	prefix := strings.ToUpper(hex.EncodeToString(prefixBytes))

	indent := strings.Repeat("    ", p.depth)
	line := fmt.Sprintf("%s%-20s # %s", indent, prefix, annotation)
	lines = append(lines, line)

	switch major {
	case 4: // Array
		for i := 0; i < length; i++ {
			if p.offset >= len(p.data) {
				lines = append(lines, fmt.Sprintf("%sERROR: truncated array", indent))
				break
			}
			lines = append(lines, p.parseItem()...)
		}
	case 5: // Map
		for i := 0; i < length*2; i++ {
			if p.offset >= len(p.data) {
				lines = append(lines, fmt.Sprintf("%sERROR: truncated map", indent))
				break
			}
			lines = append(lines, p.parseItem()...)
		}
	case 6: // Tag
		if p.offset < len(p.data) {
			lines = append(lines, p.parseItem()...)
		}
	}

	return lines
}

func (p *cborParser) parseUint(info byte) uint64 {
	if info < 24 {
		return uint64(info)
	}
	var val uint64
	start := p.offset
	switch info {
	case 24:
		if p.offset+1 > len(p.data) {
			return 0
		}
		val = uint64(p.data[p.offset])
		p.offset += 1
	case 25:
		if p.offset+2 > len(p.data) {
			return 0
		}
		val = uint64(p.data[p.offset])<<8 | uint64(p.data[p.offset+1])
		p.offset += 2
	case 26:
		if p.offset+4 > len(p.data) {
			return 0
		}
		val = uint64(p.data[p.offset])<<24 | uint64(p.data[p.offset+1])<<16 |
			uint64(p.data[p.offset+2])<<8 | uint64(p.data[p.offset+3])
		p.offset += 4
	case 27:
		if p.offset+8 > len(p.data) {
			return 0
		}
		val = uint64(p.data[p.offset])<<56 | uint64(p.data[p.offset+1])<<48 |
			uint64(p.data[p.offset+2])<<40 | uint64(p.data[p.offset+3])<<32 |
			uint64(p.data[p.offset+4])<<24 | uint64(p.data[p.offset+5])<<16 |
			uint64(p.data[p.offset+6])<<8 | uint64(p.data[p.offset+7])
		p.offset += 8
	}
	return val
}

func (p *cborParser) parseNint(info byte) int64 {
	return -1 - int64(p.parseUint(info))
}

func (p *cborParser) parseLength(info byte) int {
	if info < 24 {
		return int(info)
	}
	if p.offset >= len(p.data) {
		return 0
	}
	var length int
	switch info {
	case 24:
		length = int(p.data[p.offset])
		p.offset++
	case 25:
		length = int(p.data[p.offset])<<8 | int(p.data[p.offset+1])
		p.offset += 2
	case 26:
		length = int(p.data[p.offset])<<24 | int(p.data[p.offset+1])<<16 |
			int(p.data[p.offset+2])<<8 | int(p.data[p.offset+3])
		p.offset += 4
	case 27:
		high := uint64(p.data[p.offset])<<56 | uint64(p.data[p.offset+1])<<48 |
			uint64(p.data[p.offset+2])<<40 | uint64(p.data[p.offset+3])<<32
		low := uint64(p.data[p.offset+4])<<24 | uint64(p.data[p.offset+5])<<16 |
			uint64(p.data[p.offset+6])<<8 | uint64(p.data[p.offset+7])
		length = int(high | low)
		p.offset += 8
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
