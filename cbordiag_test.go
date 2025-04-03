package cbordiag

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestParseUint(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		info     byte
		expected uint64
		offset   int
	}{
		{"Direct value", []byte{}, 0x17, 23, 0},
		{"1-byte", []byte{0x2A}, 24, 42, 1},
		{"2-byte", []byte{0x01, 0x00}, 25, 256, 2},
		{"4-byte", []byte{0x00, 0x01, 0x00, 0x00}, 26, 65536, 4},
		{"8-byte", []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}, 27, 4294967296, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &CborParser{Data: tt.data}
			result := p.parseUint(tt.info)
			if result != tt.expected || p.Offset != tt.offset {
				t.Errorf("parseUint(%x) = %d (offset %d), want %d (offset %d)",
					tt.info, result, p.Offset, tt.expected, tt.offset)
			}
		})
	}
}

func TestParseNint(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		info     byte
		expected int64
	}{
		{"Direct value", []byte{}, 0x13, -20},
		{"1-byte", []byte{0x2A}, 24, -43},
		{"2-byte", []byte{0x01, 0x00}, 25, -257},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &CborParser{Data: tt.data}
			result := p.parseNint(tt.info)
			if result != tt.expected {
				t.Errorf("parseNint(%x) = %d, want %d", tt.info, result, tt.expected)
			}
		})
	}
}

func TestParseLength(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		info     byte
		expected int
		offset   int
	}{
		{"Direct length", []byte{}, 0x05, 5, 0},
		{"1-byte", []byte{0x20}, 24, 32, 1},
		{"2-byte", []byte{0x03, 0xE8}, 25, 1000, 2},
		{"4-byte", []byte{0x00, 0x0F, 0x42, 0x40}, 26, 1000000, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &CborParser{Data: tt.data}
			result := p.parseLength(tt.info)
			if result != tt.expected || p.Offset != tt.offset {
				t.Errorf("parseLength(%x) = %d (offset %d), want %d (offset %d)",
					tt.info, result, p.Offset, tt.expected, tt.offset)
			}
		})
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		input    []byte
		expected bool
	}{
		{[]byte("Hello"), true},
		{[]byte{0x01, 0x02}, false},
		{[]byte("123\n"), false},
		{[]byte(" café"), true}, // Contains non-ASCII but printable characters
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			result := isPrintable(tt.input)
			if result != tt.expected {
				t.Errorf("isPrintable(%x) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseItem(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []string
	}{
		{
			"Simple array",
			[]byte{0x83, 0x01, 0x02, 0x03},
			[]string{
				"83                   # ARRAY (3 items)",
				"    01                 # POS INT: 1",
				"    02                 # POS INT: 2",
				"    03                 # POS INT: 3",
			},
		},
		{
			"Nested map",
			hexDecode("a26161016162820203"),
			[]string{
				"A2                   # MAP (2 pairs)",
				"    6161               # TEXT: \"a\" (1 byte)",
				"    01                 # POS INT: 1",
				"    6162               # TEXT: \"b\" (1 byte)",
				"    82                 # ARRAY (2 items)",
				"        02              # POS INT: 2",
				"        03              # POS INT: 3",
			},
		},
		{
			"Tagged item",
			hexDecode("d82547656e65726963"),
			[]string{
				"D82547               # TAG (37)",
				"    656E65726963       # TEXT: \"Generic\" (7 bytes)",
			},
		},
		{
			"Truncated byte string",
			hexDecode("42"), // Length 2 but no data
			[]string{
				"42                   # ERROR: truncated byte string (need 2 bytes)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &CborParser{Data: tt.data, Depth: 0}
			result := p.ParseItem()
			if !compareLines(result, tt.expected) {
				t.Errorf("ParseItem() mismatch\nGot:\n%s\nWant:\n%s",
					strings.Join(result, "\n"),
					strings.Join(tt.expected, "\n"))
			}
		})
	}
}

func hexDecode(s string) []byte {
	data, _ := hex.DecodeString(s)
	return data
}

func compareLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		got := strings.TrimRight(a[i], " ")
		want := strings.TrimRight(b[i], " ")
		if got != want {
			return false
		}
	}
	return true
}
