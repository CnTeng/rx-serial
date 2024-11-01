package types

import "slices"

type Escape struct {
	EscapeChar  byte   `json:"escape_char"`
	EscapedChar []byte `json:"escaped_char"`
}

func (e *Escape) escape(data []byte) []byte {
	escaped := make([]byte, 0, len(data))
	for _, b := range data {
		if slices.Contains(e.EscapedChar, b) {
			escaped = append(escaped, byte(e.EscapeChar))
		}
		escaped = append(escaped, b)
	}

	return escaped
}

func (e *Escape) unescape(data []byte) []byte {
	unescaped := make([]byte, 0, len(data))
	for i := 0; i < len(data); i++ {
		if data[i] == byte(e.EscapeChar) {
			i++
		}
		unescaped = append(unescaped, data[i])
	}
	return unescaped
}
