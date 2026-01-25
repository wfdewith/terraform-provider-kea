package keadhcp4

import (
	"context"

	"github.com/wfdewith/terraform-provider-kea/kea"
	"github.com/wfdewith/terraform-provider-kea/kea/keaquery"
)

func (c *Client) GetSubnet(ctx context.Context, target kea.OperationTarget, query keaquery.SubnetQuery) (*Subnet, error) {
	result, err := kea.ExecWithResponse[struct {
		Subnet4 []Subnet `json:"subnet4"`
	}](ctx, c.transport, "subnet4-get", kea.WithTarget(target, query))

	if err != nil {
		return nil, err
	}

	if result == nil || len(result.Subnet4) == 0 {
		return nil, nil
	}

	return &result.Subnet4[0], nil
}

func (c *Client) AddSubnet(ctx context.Context, target kea.OperationTarget, subnet Subnet) error {
	return kea.Exec(ctx, c.transport, "subnet4-add",kea.WithTarget(target, struct {
		Subnet4 []Subnet `json:"subnet4"`
	}{
		Subnet4: []Subnet{subnet},
	}))
}

func (c *Client) UpdateSubnet(ctx context.Context, target kea.OperationTarget, subnet Subnet) error {
	return kea.Exec(ctx, c.transport, "subnet4-update", kea.WithTarget(target, struct {
		Subnet4 []Subnet `json:"subnet4"`
	}{
		Subnet4: []Subnet{subnet},
	}))
}

func (c *Client) DeleteSubnet(ctx context.Context, target kea.OperationTarget, query keaquery.SubnetQuery) error {
	return kea.Exec(ctx, c.transport, "subnet4-del", kea.WithTarget(target, query))
}
