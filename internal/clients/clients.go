package clients

import "github.com/wfdewith/terraform-provider-kea/kea/keadhcp4"

type KeaClients struct {
	DHCP4 *keadhcp4.Client
}
