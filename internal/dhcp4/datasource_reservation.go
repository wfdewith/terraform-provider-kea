package dhcp4

import (
	"context"
	"math"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/hwtypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wfdewith/terraform-provider-kea/internal/clients"
	"github.com/wfdewith/terraform-provider-kea/internal/errors"
	"github.com/wfdewith/terraform-provider-kea/kea"
)

var _ datasource.DataSourceWithConfigure = (*ReservationDataSource)(nil)
var _ datasource.DataSourceWithConfigValidators = (*ReservationDataSource)(nil)

type ReservationDataSource struct {
	client *kea.DHCP4Client
}

func (d *ReservationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp4_reservation"
}

func (d *ReservationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a DHCPv4 host reservation from the Kea DHCP server. " +
			"Host reservations bind specific DHCP resources to individual clients identified by unique identifiers. " +
			"Use this data source to look up existing reservations by client identifier.\n\n" +
			"**Important:** This data source requires the `host_cmds` hook library to be loaded on the Kea server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the reservation in the format `{subnet_id}/{identifier_type}/{identifier}`.",
			},
			"subnet_id": schema.Int64Attribute{
				Description: "The ID of the subnet containing the reservation. " +
					"Use a value of zero (0) to look up global reservations. " +
					"For subnet-specific reservations, use the subnet's numeric identifier.",
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, math.MaxUint32),
				},
			},

			"circuit_id": schema.StringAttribute{
				Description: "Circuit ID option value (Option 82 sub-option 1) to identify which reservation to retrieve. " +
					"Exactly one identifier (circuit_id, client_id, duid, flex_id, or hw_address) must be specified.",
				Optional: true,
			},
			"client_id": schema.StringAttribute{
				Description: "DHCPv4 client identifier (Option 61) to identify which reservation to retrieve. " +
					"Exactly one identifier (circuit_id, client_id, duid, flex_id, or hw_address) must be specified.",
				Optional: true,
			},
			"duid": schema.StringAttribute{
				Description: "DHCP Unique Identifier to identify which reservation to retrieve. " +
					"Exactly one identifier (circuit_id, client_id, duid, flex_id, or hw_address) must be specified.",
				Optional: true,
			},
			"flex_id": schema.StringAttribute{
				Description: "Flexible identifier to identify which reservation to retrieve. " +
					"Exactly one identifier (circuit_id, client_id, duid, flex_id, or hw_address) must be specified.",
				Optional: true,
			},
			"hw_address": schema.StringAttribute{
				Description: "Hardware (MAC) address to identify which reservation to retrieve. " +
					"Exactly one identifier (circuit_id, client_id, duid, flex_id, or hw_address) must be specified.",
				Optional:   true,
				CustomType: hwtypes.MACAddressType{},
			},

			"ip_address": schema.StringAttribute{
				Description: "IPv4 address reserved for this client.",
				Computed:    true,
				CustomType:  iptypes.IPv4AddressType{},
			},

			"boot_file_name": schema.StringAttribute{
				Description: "Boot file name (corresponds to the 'file' field in the DHCP packet and DHCP option 67) configured for this reservation.",
				Computed:    true,
			},
			"client_classes": schema.SetAttribute{
				Description: "List of client class names assigned to this reserved client.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname assigned to the reserved client.",
				Computed:    true,
			},
			"next_server": schema.StringAttribute{
				Description: "IPv4 address of the next server (corresponds to the 'siaddr' field in the DHCP packet) configured for this reservation.",
				Computed:    true,
				CustomType:  iptypes.IPv4AddressType{},
			},
			"server_hostname": schema.StringAttribute{
				Description: "Server hostname (corresponds to the 'sname' field in the DHCP packet) configured for this reservation.",
				Computed:    true,
			},
			"user_context": schema.StringAttribute{
				Description: "Arbitrary JSON data stored with this reservation.",
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
			"option_data": schema.SetNestedAttribute{
				Description: "DHCP options configured for this reservation.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the DHCP option.",
							Computed:    true,
						},
						"code": schema.Int32Attribute{
							Description: "Numeric code of the DHCP option.",
							Computed:    true,
						},
						"data": schema.StringAttribute{
							Description: "Value of the DHCP option.",
							Computed:    true,
						},
						"csv_format": schema.BoolAttribute{
							Description: "Whether the option data is in comma-separated value format (true) or hexadecimal format (false).",
							Computed:    true,
						},
						"always_send": schema.BoolAttribute{
							Description: "Whether the server always sends this option to the client, even if not requested.",
							Computed:    true,
						},
						"never_send": schema.BoolAttribute{
							Description: "Whether the server never sends this option to the client.",
							Computed:    true,
						},
						"client_classes": schema.SetAttribute{
							Description: "Client classes for which this option applies.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *ReservationDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("circuit_id"),
			path.MatchRoot("client_id"),
			path.MatchRoot("duid"),
			path.MatchRoot("flex_id"),
			path.MatchRoot("hw_address"),
		),
	}
}

func (d *ReservationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	keaClients, ok := req.ProviderData.(clients.KeaClients)
	if !ok {
		resp.Diagnostics.Append(errors.NewProviderDataTypeError(req.ProviderData))
		return
	}

	if keaClients.DHCP4 == nil {
		resp.Diagnostics.Append(errors.NewUnconfiguredClientError("dhcp4"))
		return
	}

	d.client = keaClients.DHCP4
}

func (d *ReservationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ReservationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservation, err := d.client.GetReservation(ctx, data.BuildQuery())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Reservation", err.Error())
		return
	}

	if reservation == nil {
		resp.Diagnostics.AddError("Reservation Not Found",
			"No reservation found matching the specified criteria.")
		return
	}

	resp.Diagnostics.Append(data.FromAPI(ctx, reservation)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
