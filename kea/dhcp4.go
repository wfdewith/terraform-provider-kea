package kea

import (
	"context"
	"encoding/json"
	"net/netip"
)

type DHCP4Client struct {
	transport Transport
}

func NewDHCP4Client(transport Transport) *DHCP4Client {
	return &DHCP4Client{transport}
}

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

func (c *DHCP4Client) GetSubnets(ctx context.Context) ([]Subnet, error) {
	result, err := execWithResponse[struct {
		Subnets []Subnet `json:"subnets"`
	}](ctx, c.transport, "subnet4-list", nil)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return []Subnet{}, nil
	}

	return result.Subnets, nil
}

func (c *DHCP4Client) GetReservations(ctx context.Context, subnetID uint32) ([]Reservation, error) {
	result, err := execWithResponse[struct {
		Hosts []Reservation `json:"hosts"`
	}](ctx, c.transport, "reservation-get-all", struct {
		SubnetID uint32 `json:"subnet-id"`
	}{SubnetID: subnetID})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return []Reservation{}, nil
	}

	return result.Hosts, nil
}

func (c *DHCP4Client) GetReservation(ctx context.Context, query ReservationQuery) (*Reservation, error) {
	return execWithResponse[Reservation](ctx, c.transport, "reservation-get", query)
}

func (c *DHCP4Client) AddReservation(ctx context.Context, reservation Reservation) error {
	return exec(ctx, c.transport, "reservation-add", struct {
		Reservation Reservation `json:"reservation"`
	}{
		Reservation: reservation,
	})
}

func (c *DHCP4Client) UpdateReservation(ctx context.Context, reservation Reservation) error {
	return exec(ctx, c.transport, "reservation-update", struct {
		Reservation Reservation `json:"reservation"`
	}{
		Reservation: reservation,
	})
}

func (c *DHCP4Client) DeleteReservation(ctx context.Context, query ReservationQuery) error {
	return exec(ctx, c.transport, "reservation-del", query)
}
