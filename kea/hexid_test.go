package kea_test

import (
	"testing"

	"github.com/wfdewith/terraform-provider-kea/kea"
)

func TestHexIDString(t *testing.T) {
	tests := []struct {
		name string
		id   kea.HexID
		want string
	}{
		{"nil", nil, ""},
		{"empty", kea.HexID{}, ""},
		{"single byte", kea.HexID{0x0a}, "0a"},
		{"multiple bytes", kea.HexID{0xaa, 0xbb, 0xcc}, "aa:bb:cc"},
		{"leading zero byte", kea.HexID{0x00, 0x01}, "00:01"},
		{"all byte values", kea.HexID{0x00, 0x0f, 0x10, 0xff}, "00:0f:10:ff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.want {
				t.Errorf("HexID.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
