package kea

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type HexID []byte

func ParseHexID(id string) (HexID, error) {
	if len(id) == 0 {
		return nil, nil
	}
	if strings.Contains(id, ":") {
		return parseWithSeparator(id, ":")
	}
	if strings.Contains(id, " ") {
		return parseWithSeparator(id, " ")
	}
	if strings.HasPrefix(id, "0x") {
		return parseWithPadding(id[2:])
	}
	return parseWithPadding(id)
}

func parseWithSeparator(s string, sep string) (HexID, error) {
	parts := strings.Split(s, sep)
	id := make(HexID, len(parts))
	for i, part := range parts {
		if len(part) == 1 {
			part = "0" + part
		} else if len(part) != 2 {
			return nil, fmt.Errorf("invalid segment length: %q", part)
		}

		if _, err := hex.Decode(id[i:i+1], []byte(part)); err != nil {
			return nil, err
		}
	}
	return id, nil
}

func parseWithPadding(s string) (HexID, error) {
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

func (h HexID) String() string {
	n := len(h)
	if n == 0 {
		return ""
	}

	buf := make([]byte, n*3-1)
	for i := range n {
		hex.Encode(buf[i*3:i*3+2], h[i:i+1])
		if i < n-1 {
			buf[i*3+2] = ':'
		}
	}
	return string(buf)
}

func (h HexID) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *HexID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	id, err := ParseHexID(s)
	if err != nil {
		return err
	}
	*h = id
	return nil
}
