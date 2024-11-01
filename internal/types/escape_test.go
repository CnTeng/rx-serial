package types

import (
	"testing"

	"golang.org/x/exp/slices"
)

var escape = Escape{EscapeChar: 0x10, EscapedChar: []byte{0x02, 0x03, 0x10}}

func TestEscape_escape(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x10}
	want := []byte{0x01, 0x10, 0x02, 0x10, 0x03, 0x04, 0x10, 0x10}
	result := escape.escape(data)

	if !slices.Equal(result, want) {
		t.Errorf("expected %v, got %v", want, result)
	}
}

func TestEscape_unescape(t *testing.T) {
	data := []byte{0x01, 0x10, 0x02, 0x10, 0x03, 0x04, 0x10, 0x10, 0x10, 0x10}
	want := []byte{0x01, 0x02, 0x03, 0x04, 0x10, 0x10}
	result := escape.unescape(data)

	if !slices.Equal(result, want) {
		t.Errorf("expected %v, got %v", want, result)
	}
}
