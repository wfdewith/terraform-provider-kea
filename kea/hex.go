package kea

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func ParseHex(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, nil
	}
	if strings.Contains(s, ":") {
		return parseHexWithSeparator(s, ":")
	}
	if strings.Contains(s, " ") {
		return parseHexWithSeparator(s, " ")
	}
	if strings.HasPrefix(s, "0x") {
		return parseHexWithPadding(s[2:])
	}
	return parseHexWithPadding(s)
}

func parseHexWithSeparator(s string, sep string) ([]byte, error) {
	parts := strings.Split(s, sep)
	b := make([]byte, len(parts))
	for i, part := range parts {
		if len(part) == 1 {
			part = "0" + part
		} else if len(part) != 2 {
			return nil, fmt.Errorf("invalid segment length: %q", part)
		}

		if _, err := hex.Decode(b[i:i+1], []byte(part)); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func parseHexWithPadding(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}
