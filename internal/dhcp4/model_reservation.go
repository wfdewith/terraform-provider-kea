package dhcp4

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/hwtypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wfdewith/terraform-provider-kea/kea/keadhcp4"
	"github.com/wfdewith/terraform-provider-kea/kea/keaquery"
)

type ReservationModel struct {
	ID             types.String         `tfsdk:"id"`
	SubnetID       types.Int64          `tfsdk:"subnet_id"`
	CircuitID      types.String         `tfsdk:"circuit_id"`
	ClientID       types.String         `tfsdk:"client_id"`
	DUID           types.String         `tfsdk:"duid"`
	FlexID         types.String         `tfsdk:"flex_id"`
	HWAddress      hwtypes.MACAddress   `tfsdk:"hw_address"`
	IPAddress      iptypes.IPv4Address  `tfsdk:"ip_address"`
	BootFileName   types.String         `tfsdk:"boot_file_name"`
	ClientClasses  types.Set            `tfsdk:"client_classes"`
	Hostname       types.String         `tfsdk:"hostname"`
	NextServer     iptypes.IPv4Address  `tfsdk:"next_server"`
	ServerHostname types.String         `tfsdk:"server_hostname"`
	OptionData     types.Set            `tfsdk:"option_data"`
	UserContext    jsontypes.Normalized `tfsdk:"user_context"`
}

type OptionDataModel struct {
	Name          types.String `tfsdk:"name"`
	Code          types.Int32  `tfsdk:"code"`
	Space         types.String `tfsdk:"space"`
	Data          types.String `tfsdk:"data"`
	CSVFormat     types.Bool   `tfsdk:"csv_format"`
	AlwaysSend    types.Bool   `tfsdk:"always_send"`
	NeverSend     types.Bool   `tfsdk:"never_send"`
	ClientClasses types.Set    `tfsdk:"client_classes"`
}

func (m *ReservationModel) ToAPI(ctx context.Context) (keadhcp4.Reservation, diag.Diagnostics) {
	var diags diag.Diagnostics

	reservation := keadhcp4.Reservation{
		SubnetID:       uint32(m.SubnetID.ValueInt64()),
		CircuitID:      m.CircuitID.ValueString(),
		ClientID:       m.ClientID.ValueString(),
		DUID:           m.DUID.ValueString(),
		FlexID:         m.FlexID.ValueString(),
		HWAddress:      m.HWAddress.ValueString(),
		IPAddress:      ipv4Pointer(m.IPAddress),
		BootFileName:   m.BootFileName.ValueString(),
		Hostname:       m.Hostname.ValueString(),
		NextServer:     ipv4Pointer(m.NextServer),
		ServerHostname: m.ServerHostname.ValueString(),
	}

	if !m.ClientClasses.IsNull() && !m.ClientClasses.IsUnknown() {
		diags.Append(m.ClientClasses.ElementsAs(ctx, &reservation.ClientClasses, false)...)
	}

	if !m.OptionData.IsNull() && !m.OptionData.IsUnknown() {
		var optionDataModels []OptionDataModel
		diags.Append(m.OptionData.ElementsAs(ctx, &optionDataModels, false)...)

		for _, od := range optionDataModels {
			option, d := od.ToAPI(ctx)
			diags.Append(d...)
			reservation.OptionData = append(reservation.OptionData, option)
		}
	}

	if !m.UserContext.IsNull() && !m.UserContext.IsUnknown() {
		reservation.UserContext = []byte(m.UserContext.ValueString())
	}

	return reservation, diags
}

func (m *ReservationModel) FromAPI(ctx context.Context, r *keadhcp4.Reservation) diag.Diagnostics {
	var diags diag.Diagnostics

	m.SubnetID = types.Int64Value(int64(r.SubnetID))
	m.CircuitID = stringOrNull(r.CircuitID)
	m.ClientID = stringOrNull(r.ClientID)
	m.DUID = stringOrNull(r.DUID)
	m.FlexID = stringOrNull(r.FlexID)
	m.HWAddress = macAddressOrNull(r.HWAddress)
	m.IPAddress = ipv4AddressOrNull(r.IPAddress)
	m.BootFileName = stringOrNull(r.BootFileName)
	m.Hostname = stringOrNull(r.Hostname)
	m.NextServer = ipv4AddressOrNull(r.NextServer)
	m.ServerHostname = stringOrNull(r.ServerHostname)

	if len(r.ClientClasses) > 0 {
		setVal, d := types.SetValueFrom(ctx, types.StringType, r.ClientClasses)
		diags.Append(d...)
		m.ClientClasses = setVal
	} else {
		m.ClientClasses = types.SetNull(types.StringType)
	}

	if len(r.OptionData) > 0 {
		var optionDataModels []OptionDataModel
		for _, od := range r.OptionData {
			var odModel OptionDataModel
			diags.Append(odModel.FromAPI(ctx, &od)...)
			optionDataModels = append(optionDataModels, odModel)
		}
		optionDataSet, d := types.SetValueFrom(ctx, m.OptionData.ElementType(ctx), optionDataModels)
		diags.Append(d...)
		m.OptionData = optionDataSet
	} else {
		m.OptionData = types.SetNull(m.OptionData.ElementType(ctx))
	}

	if len(r.UserContext) > 0 {
		m.UserContext = jsontypes.NewNormalizedValue(string(r.UserContext))
	} else {
		m.UserContext = jsontypes.NewNormalizedNull()
	}

	m.ID = types.StringValue(m.ComputeID())

	return diags
}

func (m *ReservationModel) BuildQuery() keaquery.ReservationQuery {
	subnetID := uint32(m.SubnetID.ValueInt64())

	identifierType, identifier := m.getIdentifier()
	if identifierType != "" {
		return keaquery.ReservationByIdentifier(subnetID, identifierType, identifier)
	}

	return keaquery.ReservationByIP(subnetID, m.IPAddress.ValueString())
}

func (m *ReservationModel) ComputeID() string {
	subnetID := m.SubnetID.ValueInt64()

	identifierType, identifier := m.getIdentifier()
	if identifierType != "" {
		return fmt.Sprintf("%d/%s/%s", subnetID, identifierType, identifier)
	}

	return fmt.Sprintf("%d/ip-address/%s", subnetID, m.IPAddress.ValueString())
}

func (m *ReservationModel) getIdentifier() (identifierType, identifier string) {
	switch {
	case !m.CircuitID.IsNull() && !m.CircuitID.IsUnknown():
		return "circuit-id", m.CircuitID.ValueString()
	case !m.ClientID.IsNull() && !m.ClientID.IsUnknown():
		return "client-id", m.ClientID.ValueString()
	case !m.DUID.IsNull() && !m.DUID.IsUnknown():
		return "duid", m.DUID.ValueString()
	case !m.FlexID.IsNull() && !m.FlexID.IsUnknown():
		return "flex-id", m.FlexID.ValueString()
	case !m.HWAddress.IsNull() && !m.HWAddress.IsUnknown():
		return "hw-address", m.HWAddress.ValueString()
	default:
		return "", ""
	}
}

func (o *OptionDataModel) ToAPI(ctx context.Context) (keadhcp4.OptionData, diag.Diagnostics) {
	var diags diag.Diagnostics

	option := keadhcp4.OptionData{
		Name:       o.Name.ValueString(),
		Code:       uint8(o.Code.ValueInt32()),
		Space:      o.Space.ValueString(),
		Data:       o.Data.ValueString(),
		CSVFormat:  boolPointer(o.CSVFormat),
		AlwaysSend: boolPointer(o.AlwaysSend),
		NeverSend:  boolPointer(o.NeverSend),
	}

	if !o.ClientClasses.IsNull() && !o.ClientClasses.IsUnknown() {
		diags.Append(o.ClientClasses.ElementsAs(ctx, &option.ClientClasses, false)...)
	}

	return option, diags
}

func (o *OptionDataModel) FromAPI(ctx context.Context, od *keadhcp4.OptionData) diag.Diagnostics {
	var diags diag.Diagnostics

	o.Name = types.StringValue(od.Name)
	o.Code = types.Int32Value(int32(od.Code))
	o.Space = stringOrNull(od.Space)
	o.Data = stringOrNull(od.Data)
	o.CSVFormat = types.BoolPointerValue(od.CSVFormat)
	o.AlwaysSend = types.BoolPointerValue(od.AlwaysSend)
	o.NeverSend = types.BoolPointerValue(od.NeverSend)

	if len(od.ClientClasses) > 0 {
		setVal, d := types.SetValueFrom(ctx, types.StringType, od.ClientClasses)
		diags.Append(d...)
		o.ClientClasses = setVal
	} else {
		o.ClientClasses = types.SetNull(types.StringType)
	}

	return diags
}

func stringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func macAddressOrNull(s string) hwtypes.MACAddress {
	if s == "" {
		return hwtypes.NewMACAddressNull()
	}
	return hwtypes.NewMACAddressValue(s)
}

func ipv4AddressOrNull(ip *netip.Addr) iptypes.IPv4Address {
	if ip == nil || !ip.IsValid() || ip.IsUnspecified() {
		return iptypes.NewIPv4AddressNull()
	}
	return iptypes.NewIPv4AddressValue(ip.String())
}

func boolPointer(b types.Bool) *bool {
	if b.IsNull() || b.IsUnknown() {
		return nil
	}

	r := b.ValueBool()
	return &r
}

func ipv4Pointer(ip iptypes.IPv4Address) *netip.Addr {
	if ip.IsNull() || ip.IsUnknown() {
		return nil
	}

	r, _ := ip.ValueIPv4Address()
	return &r
}
