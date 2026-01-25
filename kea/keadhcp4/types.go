package keadhcp4

import (
	"encoding/json"
	"net/netip"
)

type Subnet struct {
	ID     uint32 `json:"id"`
	Subnet string `json:"subnet"`

	SharedNetworkName string `json:"shared-network-name,omitempty"`
}

type Reservation struct {
	SubnetID uint32 `json:"subnet-id"`

	CircuitID string `json:"circuit-id,omitempty"`
	ClientID  string `json:"client-id,omitempty"`
	DUID      string `json:"duid,omitempty"`
	FlexID    string `json:"flex-id,omitempty"`
	HWAddress string `json:"hw-address,omitempty"`

	BootFileName   string          `json:"boot-file-name,omitempty"`
	ClientClasses  []string        `json:"client-classes,omitempty"`
	Hostname       string          `json:"hostname,omitempty"`
	IPAddress      *netip.Addr     `json:"ip-address,omitempty"`
	NextServer     *netip.Addr     `json:"next-server,omitempty"`
	OptionData     []OptionData    `json:"option-data,omitempty"`
	ServerHostname string          `json:"server-hostname,omitempty"`
	UserContext    json.RawMessage `json:"user-context,omitempty"`
}

type OptionData struct {
	Name          string   `json:"name,omitempty"`
	Data          string   `json:"data,omitempty"`
	Code          uint8    `json:"code,omitempty"`
	Space         string   `json:"space,omitempty"`
	CSVFormat     *bool    `json:"csv-format,omitempty"`
	AlwaysSend    *bool    `json:"always-send,omitempty"`
	NeverSend     *bool    `json:"never-send,omitempty"`
	ClientClasses []string `json:"client-classes,omitempty"`
}
