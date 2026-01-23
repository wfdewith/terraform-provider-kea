package dhcp4

import (
	"context"
	"math"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/hwtypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wfdewith/terraform-provider-kea/internal/clients"
	"github.com/wfdewith/terraform-provider-kea/internal/errors"
	"github.com/wfdewith/terraform-provider-kea/kea"
)

var _ resource.ResourceWithConfigure = (*ReservationResource)(nil)
var _ resource.ResourceWithConfigValidators = (*ReservationResource)(nil)

type ReservationResource struct {
	client *kea.DHCP4Client
}

func (r *ReservationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp4_reservation"
}

func (r *ReservationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DHCPv4 host reservation in the Kea DHCP server. " +
			"Host reservations allow binding specific DHCP resources (such as IP addresses and options) to individual clients " +
			"identified by unique identifiers like MAC addresses, client IDs, or circuit IDs. " +
			"Reservations can be scoped to a specific subnet or configured globally (subnet_id = 0).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the reservation in the format `{subnet_id}/{identifier_type}/{identifier}`.",
			},
			"subnet_id": schema.Int64Attribute{
				Description: "The ID of the subnet to which the reservation belongs. " +
					"Use a value of zero (0) to create a global reservation that applies across all subnets. " +
					"For subnet-specific reservations, use the subnet's numeric identifier as defined in the Kea configuration.",
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, math.MaxUint32),
				},
			},

			"circuit_id": schema.StringAttribute{
				Description: "Circuit ID option value (Option 82 sub-option 1) used to identify the DHCP client. " +
					"Typically inserted by DHCP relay agents. " +
					"Only one identifier type (circuit_id, client_id, duid, flex_id, or hw_address) should be specified per reservation.",
				Optional: true,
			},
			"client_id": schema.StringAttribute{
				Description: "DHCPv4 client identifier (Option 61) used to identify the DHCP client. " +
					"This is a unique identifier sent by the client in DHCP messages. " +
					"Only one identifier type (circuit_id, client_id, duid, flex_id, or hw_address) should be specified per reservation.",
				Optional: true,
			},
			"duid": schema.StringAttribute{
				Description: "DHCP Unique Identifier used to identify the client. " +
					"While typically used in DHCPv6, it can also be used for DHCPv4 client identification. " +
					"Only one identifier type (circuit_id, client_id, duid, flex_id, or hw_address) should be specified per reservation.",
				Optional: true,
			},
			"flex_id": schema.StringAttribute{
				Description: "Flexible identifier used to identify the client. " +
					"This is a custom identifier configured via the flex_id hook library, allowing flexible client identification based on various packet attributes. " +
					"Only one identifier type (circuit_id, client_id, duid, flex_id, or hw_address) should be specified per reservation.",
				Optional: true,
			},
			"hw_address": schema.StringAttribute{
				Description: "Hardware (MAC) address of the client's network interface. " +
					"This is the most commonly used identifier for DHCPv4 reservations. " +
					"Only one identifier type (circuit_id, client_id, duid, flex_id, or hw_address) should be specified per reservation.",
				Optional:   true,
				CustomType: hwtypes.MACAddressType{},
			},

			"ip_address": schema.StringAttribute{
				Description: "IPv4 address to reserve for this client. " +
					"When specified, the Kea server will always assign this address to the identified client. " +
					"The address must be within the range of the subnet (if subnet_id is non-zero) or can be any valid IPv4 address for global reservations.",
				Optional:   true,
				CustomType: iptypes.IPv4AddressType{},
			},

			"boot_file_name": schema.StringAttribute{
				Description: "Boot file name (corresponds to the 'file' field in the DHCP packet and DHCP option 67). " +
					"Used for network booting scenarios to specify the boot file that the client should download.",
				Optional: true,
			},
			"client_classes": schema.SetAttribute{
				Description: "List of client class names to assign to this reserved client. " +
					"Client classes allow different groups of clients to receive different DHCP configuration. " +
					"Classes must be defined in the Kea server configuration.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname to assign to the reserved client. " +
					"This can be used for dynamic DNS updates and to populate the hostname option (Option 12) sent to the client.",
				Optional: true,
			},
			"next_server": schema.StringAttribute{
				Description: "IPv4 address of the next server to use in the boot process (corresponds to the 'siaddr' field in the DHCP packet). " +
					"Typically used in network boot scenarios to specify the TFTP server address.",
				Optional:   true,
				Computed:   true,
				CustomType: iptypes.IPv4AddressType{},
			},
			"server_hostname": schema.StringAttribute{
				Description: "Server hostname (corresponds to the 'sname' field in the DHCP packet). " +
					"Used in network boot scenarios to specify the name of the boot server.",
				Optional: true,
			},
			"user_context": schema.StringAttribute{
				Description: "Arbitrary JSON data to store custom information with this reservation. " +
					"This can be used to store additional metadata such as contact information, asset tags, or other operational data. " +
					"The data is stored but not processed by Kea.",
				Optional:   true,
				Computed:   true,
				CustomType: jsontypes.NormalizedType{},
			},
		},
		Blocks: map[string]schema.Block{
			"option_data": schema.SetNestedBlock{
				Description: "DHCP options to be sent to the reserved client. " +
					"These options override subnet-level and global options for this specific client. " +
					"Options can be specified by name or numeric code.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the DHCP option (e.g., 'domain-name-servers', 'routers'). " +
								"Either name or code must be specified. If both are specified, they must refer to the same DHCP option " +
								"(e.g., name='routers' and code=3).",
							Optional: true,
							Computed: true,
						},
						"code": schema.Int32Attribute{
							Description: "Numeric code of the DHCP option (0-255). " +
								"Either name or code must be specified. If both are specified, they must refer to the same DHCP option " +
								"(e.g., name='routers' and code=3).",
							Optional: true,
							Computed: true,
							Validators: []validator.Int32{
								int32validator.Between(0, math.MaxUint8),
							},
						},
						"space": schema.StringAttribute{
							Description: "Option space name. " +
								"Defaults to 'dhcp4' for DHCPv4 options. " +
								"Use custom space names for vendor-specific or custom option spaces.",
							Optional: true,
							Computed: true,
						},
						"data": schema.StringAttribute{
							Description: "Value of the DHCP option. " +
								"Format depends on the option type and csv_format setting. " +
								"When csv_format is true (default), data is comma-separated values (e.g., '192.0.2.1,192.0.2.2' for DNS servers). " +
								"When csv_format is false, data is a hexadecimal string (e.g., 'c000020a' for IP 192.0.2.10).",
							Optional: true,
						},
						"csv_format": schema.BoolAttribute{
							Description: "Format of the option data. " +
								"When true (default), data is interpreted as comma-separated values appropriate for the option type. " +
								"When false, data is interpreted as a raw hexadecimal string.",
							Optional: true,
							Computed: true,
						},
						"always_send": schema.BoolAttribute{
							Description: "When true, the server sends this option to the client even if the client did not request it in the Parameter Request List (Option 55). " +
								"Useful for ensuring critical options are always delivered.",
							Optional: true,
							Computed: true,
						},
						"never_send": schema.BoolAttribute{
							Description: "When true, prevents the server from sending this option to the client, even if it would normally be sent. " +
								"Useful for explicitly blocking certain options for specific clients.",
							Optional: true,
							Computed: true,
						},
						"client_classes": schema.SetAttribute{
							Description: "List of client class names for which this option applies. " +
								"If specified, the option is only sent to clients that are members of at least one of these classes.",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
					Validators: []validator.Object{
						objectvalidator.AtLeastOneOf(
							path.MatchRoot("name"),
							path.MatchRoot("code"),
						),
					},
				},
			},
		},
	}
}

func (r *ReservationResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("circuit_id"),
			path.MatchRoot("client_id"),
			path.MatchRoot("duid"),
			path.MatchRoot("flex_id"),
			path.MatchRoot("hw_address"),
		),
	}
}

func (r *ReservationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = keaClients.DHCP4
}

func (r *ReservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ReservationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservation, diags := data.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AddReservation(ctx, reservation); err != nil {
		resp.Diagnostics.AddError("Error Adding Reservation", err.Error())
		return
	}

	created, err := r.client.GetReservation(ctx, data.BuildQuery())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Created Reservation", err.Error())
		return
	}

	if created == nil {
		resp.Diagnostics.AddError("Error Reading Created Reservation", "Reservation not found after creation")
		return
	}

	resp.Diagnostics.Append(data.FromAPI(ctx, created)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ReservationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservation, err := r.client.GetReservation(ctx, data.BuildQuery())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Reservation", err.Error())
		return
	}

	if reservation == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.FromAPI(ctx, reservation)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ReservationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservation, diags := data.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateReservation(ctx, reservation); err != nil {
		resp.Diagnostics.AddError("Error Updating Reservation", err.Error())
		return
	}

	updated, err := r.client.GetReservation(ctx, data.BuildQuery())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Updated Reservation", err.Error())
		return
	}

	if updated == nil {
		resp.Diagnostics.AddError("Error Reading Updated Reservation", "Reservation not found after update")
		return
	}

	resp.Diagnostics.Append(data.FromAPI(ctx, updated)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ReservationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteReservation(ctx, data.BuildQuery()); err != nil {
		resp.Diagnostics.AddError("Error Deleting Reservation", err.Error())
	}
}
