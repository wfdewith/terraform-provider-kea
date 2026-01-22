package clients

import "github.com/wfdewith/terraform-provider-kea/kea"

type KeaClients struct {
	DHCP4 *kea.DHCP4Client
}
