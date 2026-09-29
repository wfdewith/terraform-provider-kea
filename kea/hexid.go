package kea

import (
	"encoding/json"
	"strings"
)

type HexID []byte

func ParseHexID(id string) (HexID, error) {
	b, err := ParseHex(id)
	if err != nil {
		return nil, err
	}
	return HexID(b), nil
}

func (h HexID) String() string {
	if len(h) == 0 {
		return ""
	}

	const hexDigits = "0123456789abcdef"
	var sb strings.Builder
	sb.Grow(len(h)*3 - 1)

	for i, b := range h {
		if i > 0 {
			sb.WriteByte(':')
		}
		sb.WriteByte(hexDigits[b>>4])
		sb.WriteByte(hexDigits[b&0x0f])
	}
	return sb.String()
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
