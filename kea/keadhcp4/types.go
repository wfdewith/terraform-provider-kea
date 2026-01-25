package keadhcp4

import (
	"encoding/json"
	"net/netip"

	"github.com/wfdewith/terraform-provider-kea/kea"
)

type Subnet struct {
	ID     uint32       `json:"id"`
	Subnet netip.Prefix `json:"subnet"`

	// Basic timing parameters
	ValidLifetime    *uint32 `json:"valid-lifetime,omitempty"`
	MinValidLifetime *uint32 `json:"min-valid-lifetime,omitempty"`
	MaxValidLifetime *uint32 `json:"max-valid-lifetime,omitempty"`
	RenewTimer       *uint32 `json:"renew-timer,omitempty"`
	RebindTimer      *uint32 `json:"rebind-timer,omitempty"`

	// Interface binding
	Interface *string `json:"interface,omitempty"`

	// Client identification
	MatchClientID *bool `json:"match-client-id,omitempty"`
	Authoritative *bool `json:"authoritative,omitempty"`

	// Boot parameters
	NextServer     *netip.Addr `json:"next-server,omitempty"`
	ServerHostname *string     `json:"server-hostname,omitempty"`
	BootFileName   *string     `json:"boot-file-name,omitempty"`

	// Client classification
	ClientClass               *string  `json:"client-class,omitempty"`
	ClientClasses             []string `json:"client-classes,omitempty"`
	RequireClientClasses      []string `json:"require-client-classes,omitempty"`
	EvaluateAdditionalClasses []string `json:"evaluate-additional-classes,omitempty"`

	// Reservation behavior
	ReservationsGlobal    *bool `json:"reservations-global,omitempty"`
	ReservationsInSubnet  *bool `json:"reservations-in-subnet,omitempty"`
	ReservationsOutOfPool *bool `json:"reservations-out-of-pool,omitempty"`

	// Tee times (T1/T2 calculation)
	CalculateTeeTimes *bool    `json:"calculate-tee-times,omitempty"`
	T1Percent         *float64 `json:"t1-percent,omitempty"`
	T2Percent         *float64 `json:"t2-percent,omitempty"`

	// Cache parameters
	CacheThreshold *float64 `json:"cache-threshold,omitempty"`
	CacheMaxAge    *uint32  `json:"cache-max-age,omitempty"`

	// DDNS parameters
	DDNSSendUpdates            *bool    `json:"ddns-send-updates,omitempty"`
	DDNSOverrideNoUpdate       *bool    `json:"ddns-override-no-update,omitempty"`
	DDNSOverrideClientUpdate   *bool    `json:"ddns-override-client-update,omitempty"`
	DDNSReplaceClientName      *string  `json:"ddns-replace-client-name,omitempty"`
	DDNSGeneratedPrefix        *string  `json:"ddns-generated-prefix,omitempty"`
	DDNSQualifyingSuffix       *string  `json:"ddns-qualifying-suffix,omitempty"`
	DDNSUpdateOnRenew          *bool    `json:"ddns-update-on-renew,omitempty"`
	DDNSUseConflictResolution  *bool    `json:"ddns-use-conflict-resolution,omitempty"`
	DDNSConflictResolutionMode *string  `json:"ddns-conflict-resolution-mode,omitempty"`
	DDNSTTLPercent             *float64 `json:"ddns-ttl-percent,omitempty"`
	DDNSTTL                    *uint32  `json:"ddns-ttl,omitempty"`
	DDNSTTLMin                 *uint32  `json:"ddns-ttl-min,omitempty"`
	DDNSTTLMax                 *uint32  `json:"ddns-ttl-max,omitempty"`

	// Hostname sanitization
	HostnameCharSet         *string `json:"hostname-char-set,omitempty"`
	HostnameCharReplacement *string `json:"hostname-char-replacement,omitempty"`

	// Extended info storage
	StoreExtendedInfo *bool `json:"store-extended-info,omitempty"`

	// Allocator
	Allocator *string `json:"allocator,omitempty"`

	// Offer lifetime
	OfferLifetime *uint32 `json:"offer-lifetime,omitempty"`

	// 4o6 (DHCPv4-over-DHCPv6) parameters
	Subnet4o6Interface   *string `json:"4o6-interface,omitempty"`
	Subnet4o6InterfaceID *string `json:"4o6-interface-id,omitempty"`
	Subnet4o6Subnet      *string `json:"4o6-subnet,omitempty"`

	// Relay configuration
	Relay *Relay `json:"relay,omitempty"`

	// Pools
	Pools []Pool `json:"pools,omitempty"`

	// Options
	OptionData []OptionData `json:"option-data,omitempty"`

	// User context and comment
	UserContext json.RawMessage `json:"user-context,omitempty"`
	Comment     *string         `json:"comment,omitempty"`
}

type Relay struct {
	IPAddresses []string `json:"ip-addresses,omitempty"`
}

type Pool struct {
	Pool   *string `json:"pool,omitempty"`
	PoolID *uint32 `json:"pool-id,omitempty"`

	// Client classification
	ClientClass               *string  `json:"client-class,omitempty"`
	ClientClasses             []string `json:"client-classes,omitempty"`
	RequireClientClasses      []string `json:"require-client-classes,omitempty"`
	EvaluateAdditionalClasses []string `json:"evaluate-additional-classes,omitempty"`

	// DDNS parameters (pool-level overrides)
	DDNSSendUpdates            *bool    `json:"ddns-send-updates,omitempty"`
	DDNSOverrideNoUpdate       *bool    `json:"ddns-override-no-update,omitempty"`
	DDNSOverrideClientUpdate   *bool    `json:"ddns-override-client-update,omitempty"`
	DDNSReplaceClientName      *string  `json:"ddns-replace-client-name,omitempty"`
	DDNSGeneratedPrefix        *string  `json:"ddns-generated-prefix,omitempty"`
	DDNSQualifyingSuffix       *string  `json:"ddns-qualifying-suffix,omitempty"`
	DDNSUpdateOnRenew          *bool    `json:"ddns-update-on-renew,omitempty"`
	DDNSConflictResolutionMode *string  `json:"ddns-conflict-resolution-mode,omitempty"`
	DDNSTTLPercent             *float64 `json:"ddns-ttl-percent,omitempty"`
	DDNSTTL                    *uint32  `json:"ddns-ttl,omitempty"`
	DDNSTTLMin                 *uint32  `json:"ddns-ttl-min,omitempty"`
	DDNSTTLMax                 *uint32  `json:"ddns-ttl-max,omitempty"`

	// Hostname sanitization
	HostnameCharSet         *string `json:"hostname-char-set,omitempty"`
	HostnameCharReplacement *string `json:"hostname-char-replacement,omitempty"`

	// Options
	OptionData []OptionData `json:"option-data,omitempty"`

	// User context and comment
	UserContext json.RawMessage `json:"user-context,omitempty"`
	Comment     *string         `json:"comment,omitempty"`
}

type Reservation struct {
	SubnetID uint32 `json:"subnet-id"`

	CircuitID kea.HexID `json:"circuit-id,omitempty"`
	ClientID  kea.HexID `json:"client-id,omitempty"`
	DUID      kea.HexID `json:"duid,omitempty"`
	FlexID    kea.HexID `json:"flex-id,omitempty"`
	HWAddress kea.HexID `json:"hw-address,omitempty"`

	BootFileName   *string          `json:"boot-file-name,omitempty"`
	ClientClasses  []string         `json:"client-classes,omitempty"`
	Hostname       *string          `json:"hostname,omitempty"`
	IPAddress      *netip.Addr      `json:"ip-address,omitempty"`
	NextServer     *netip.Addr      `json:"next-server,omitempty"`
	OptionData     []OptionData     `json:"option-data,omitempty"`
	ServerHostname *string          `json:"server-hostname,omitempty"`
	UserContext    *json.RawMessage `json:"user-context,omitempty"`
}

type OptionData struct {
	Name          *string  `json:"name,omitempty"`
	Data          *string  `json:"data,omitempty"`
	Code          *uint8   `json:"code,omitempty"`
	Space         *string  `json:"space,omitempty"`
	CSVFormat     *bool    `json:"csv-format,omitempty"`
	AlwaysSend    *bool    `json:"always-send,omitempty"`
	NeverSend     *bool    `json:"never-send,omitempty"`
	ClientClasses []string `json:"client-classes,omitempty"`
}

func (s *Subnet) UnmarshalJSON(data []byte) error {
	type Alias Subnet
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	s.Interface = nilIfEmpty(s.Interface)
	s.ServerHostname = nilIfEmpty(s.ServerHostname)
	s.BootFileName = nilIfEmpty(s.BootFileName)
	s.ClientClass = nilIfEmpty(s.ClientClass)
	s.DDNSReplaceClientName = nilIfEmpty(s.DDNSReplaceClientName)
	s.DDNSGeneratedPrefix = nilIfEmpty(s.DDNSGeneratedPrefix)
	s.DDNSQualifyingSuffix = nilIfEmpty(s.DDNSQualifyingSuffix)
	s.DDNSConflictResolutionMode = nilIfEmpty(s.DDNSConflictResolutionMode)
	s.HostnameCharSet = nilIfEmpty(s.HostnameCharSet)
	s.HostnameCharReplacement = nilIfEmpty(s.HostnameCharReplacement)
	s.Allocator = nilIfEmpty(s.Allocator)
	s.Subnet4o6Interface = nilIfEmpty(s.Subnet4o6Interface)
	s.Subnet4o6InterfaceID = nilIfEmpty(s.Subnet4o6InterfaceID)
	s.Subnet4o6Subnet = nilIfEmpty(s.Subnet4o6Subnet)
	s.Comment = nilIfEmpty(s.Comment)

	if s.Relay != nil && len(s.Relay.IPAddresses) == 0 {
		s.Relay = nil
	}

	if len(s.ClientClasses) == 0 {
		s.ClientClasses = nil
	}
	if len(s.RequireClientClasses) == 0 {
		s.RequireClientClasses = nil
	}
	if len(s.EvaluateAdditionalClasses) == 0 {
		s.EvaluateAdditionalClasses = nil
	}
	if len(s.Pools) == 0 {
		s.Pools = nil
	}
	if len(s.OptionData) == 0 {
		s.OptionData = nil
	}

	return nil
}

func (p *Pool) UnmarshalJSON(data []byte) error {
	type Alias Pool
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	p.Pool = nilIfEmpty(p.Pool)
	p.ClientClass = nilIfEmpty(p.ClientClass)
	p.DDNSReplaceClientName = nilIfEmpty(p.DDNSReplaceClientName)
	p.DDNSGeneratedPrefix = nilIfEmpty(p.DDNSGeneratedPrefix)
	p.DDNSQualifyingSuffix = nilIfEmpty(p.DDNSQualifyingSuffix)
	p.DDNSConflictResolutionMode = nilIfEmpty(p.DDNSConflictResolutionMode)
	p.HostnameCharSet = nilIfEmpty(p.HostnameCharSet)
	p.HostnameCharReplacement = nilIfEmpty(p.HostnameCharReplacement)
	p.Comment = nilIfEmpty(p.Comment)

	if len(p.ClientClasses) == 0 {
		p.ClientClasses = nil
	}
	if len(p.RequireClientClasses) == 0 {
		p.RequireClientClasses = nil
	}
	if len(p.EvaluateAdditionalClasses) == 0 {
		p.EvaluateAdditionalClasses = nil
	}
	if len(p.OptionData) == 0 {
		p.OptionData = nil
	}

	return nil
}

func (o *OptionData) UnmarshalJSON(data []byte) error {
	type Alias OptionData
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if len(o.ClientClasses) == 0 {
		o.ClientClasses = nil
	}

	return nil
}

func (r *Reservation) UnmarshalJSON(data []byte) error {
	type Alias Reservation
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// The Kea API returns these fields as empty strings when not set,
	// rather than omitting them from the response.
	r.BootFileName = nilIfEmpty(r.BootFileName)
	r.Hostname = nilIfEmpty(r.Hostname)
	r.ServerHostname = nilIfEmpty(r.ServerHostname)

	// The Kea API returns "0.0.0.0" for next-server when not set.
	if r.NextServer != nil && (!r.NextServer.IsValid() || r.NextServer.IsUnspecified()) {
		r.NextServer = nil
	}

	return nil
}

func nilIfEmpty(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
