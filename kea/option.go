package kea

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

// OptionPayload is the value carried by an option. The concrete type decides
// how Kea interprets it, so csv-format is derived rather than stored.
type OptionPayload interface {
	isOptionPayload()
}

type CSVPayload string

type BinaryPayload []byte

// UnresolvedPayload is data sent without csv-format. Kea reads it as CSV when
// the option has a definition and as binary otherwise, and reports which in
// its response, so it only ever appears in requests.
type UnresolvedPayload string

func (CSVPayload) isOptionPayload()        {}
func (BinaryPayload) isOptionPayload()     {}
func (UnresolvedPayload) isOptionPayload() {}

// CSVFormat is the tri-state of Kea's csv-format field. Unspecified means the
// field was omitted and Kea decides from the option definition.
type CSVFormat int

const (
	CSVFormatUnspecified CSVFormat = iota
	CSVFormatCSV
	CSVFormatBinary
)

// CSVFormatOf maps a wire csv-format value, with nil meaning omitted.
func CSVFormatOf(v *bool) CSVFormat {
	switch {
	case v == nil:
		return CSVFormatUnspecified
	case *v:
		return CSVFormatCSV
	default:
		return CSVFormatBinary
	}
}

// PayloadDescribes reports whether Kea storing want results in have.
func PayloadDescribes(want, have OptionPayload) bool {
	switch w := want.(type) {
	case nil:
		return have == nil
	case CSVPayload:
		h, ok := have.(CSVPayload)
		return ok && w == h
	case BinaryPayload:
		h, ok := have.(BinaryPayload)
		return ok && bytes.Equal(w, h)
	case UnresolvedPayload:
		switch h := have.(type) {
		case CSVPayload:
			return string(w) == string(h)
		case BinaryPayload:
			b, err := ParseBinaryPayload(string(w))
			return err == nil && bytes.Equal(b, h)
		}
	}
	return false
}

// NewOptionPayload interprets option data the way Kea does for a given
// csv-format. Unspecified means csv-format was omitted: Kea reads the data as
// CSV when the option has a definition and as binary otherwise, so the result
// stays unresolved until Kea reports the format it chose.
func NewOptionPayload(data string, format CSVFormat) (OptionPayload, error) {
	switch format {
	case CSVFormatUnspecified:
		return UnresolvedPayload(data), nil
	case CSVFormatCSV:
		return CSVPayload(data), nil
	case CSVFormatBinary:
		b, err := ParseBinaryPayload(data)
		if err != nil {
			return nil, fmt.Errorf("option data %q: %w", data, err)
		}
		return b, nil
	default:
		return nil, fmt.Errorf("unsupported csv format %d", format)
	}
}

// String renders the payload the way Kea echoes binary option data back.
func (b BinaryPayload) String() string {
	return strings.ToUpper(hex.EncodeToString(b))
}

// ParseBinaryPayload mirrors Kea's OptionDataParser for csv-format=false:
// quotedStringToBinary first, then decodeFormattedHexString on the raw input.
func ParseBinaryPayload(s string) (BinaryPayload, error) {
	if b := quotedStringToBinary(s); len(b) > 0 {
		return b, nil
	}
	if s == "" {
		return BinaryPayload{}, nil
	}
	b, err := decodeFormattedHex(s)
	if err != nil {
		return nil, fmt.Errorf("not a quoted string or a string of hexadecimal digits: %w", err)
	}
	return b, nil
}

func quotedStringToBinary(s string) []byte {
	t := strings.TrimSpace(s)
	if len(t) > 1 && t[0] == '\'' && t[len(t)-1] == '\'' {
		return []byte(strings.TrimSpace(t[1 : len(t)-1]))
	}
	return nil
}

// Unlike ParseHex, Kea rejects a bare "0x".
func decodeFormattedHex(s string) ([]byte, error) {
	if strings.Contains(s, ":") || strings.Contains(s, " ") {
		return ParseHex(s)
	}
	if len(s) > 2 && strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}
