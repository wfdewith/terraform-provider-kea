package keadhcp4

import (
	"context"

	"github.com/wfdewith/terraform-provider-kea/kea"
	"github.com/wfdewith/terraform-provider-kea/kea/keaquery"
)

func (c *Client) GetSubnet(ctx context.Context, query keaquery.SubnetQuery) (*Subnet, error) {
	result, err := kea.ExecWithResponse[struct {
		Subnet4 []Subnet `json:"subnet4"`
	}](ctx, c.transport, "subnet4-get", query)

	if err != nil {
		return nil, err
	}

	if result == nil || len(result.Subnet4) == 0 {
		return nil, nil
	}

	return &result.Subnet4[0], nil
}

func (c *Client) AddSubnet(ctx context.Context, subnet Subnet) error {
	return kea.Exec(ctx, c.transport, "subnet4-add", struct {
		Subnet4 []Subnet `json:"subnet4"`
	}{
		Subnet4: []Subnet{subnet},
	})
}

func (c *Client) UpdateSubnet(ctx context.Context, subnet Subnet) error {
	return kea.Exec(ctx, c.transport, "subnet4-update", struct {
		Subnet4 []Subnet `json:"subnet4"`
	}{
		Subnet4: []Subnet{subnet},
	})
}

func (c *Client) DeleteSubnet(ctx context.Context, query keaquery.SubnetQuery) error {
	return kea.Exec(ctx, c.transport, "subnet4-del", query)
}
