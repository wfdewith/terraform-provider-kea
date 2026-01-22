package acctest

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/wfdewith/terraform-provider-kea/internal/provider"
	"github.com/wfdewith/terraform-provider-kea/kea"
)

// ProtoV6ProviderFactories is used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kea": providerserver.NewProtocol6WithError(provider.New("test")()),
}

var (
	subnetID     uint32
	subnetPrefix string
	setupOnce    sync.Once
	setupErr     error
)

// PreCheck validates required environment variables and queries the
// Kea server for available subnets before running acceptance tests.
func PreCheck(t *testing.T) {
	t.Helper()

	address := os.Getenv("KEA_DHCP4_ADDRESS")
	if address == "" {
		t.Fatal("KEA_DHCP4_ADDRESS must be set for acceptance tests")
	}

	// Run setup only once across all tests
	setupOnce.Do(func() {
		setupErr = discoverSubnets(address)
	})

	if setupErr != nil {
		t.Fatalf("Failed to discover subnets: %s", setupErr)
	}

	if subnetID == 0 {
		t.Fatal("No subnets configured on Kea server. At least one subnet is required for acceptance tests.")
	}
}

func discoverSubnets(address string) error {
	transport := &kea.HTTPTransport{Endpoint: address}
	client := kea.NewDHCP4Client(transport)

	subnets, err := client.GetSubnets(context.Background())
	if err != nil {
		return err
	}

	if len(subnets) > 0 {
		subnetID = subnets[0].ID
		subnetPrefix = subnets[0].Subnet
	}

	return nil
}

func ProviderConfig() string {
	return `
provider "kea" {}
`
}

// SubnetID returns a valid subnet ID discovered from the Kea server.
func SubnetID() uint32 {
	return subnetID
}

// SubnetPrefix returns the subnet prefix (CIDR) discovered from the Kea server.
func SubnetPrefix() string {
	return subnetPrefix
}

// RandomMAC generates a random locally administered MAC address.
func RandomMAC() string {
	return fmt.Sprintf("02:%02x:%02x:%02x:%02x:%02x",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
	)
}

// RandomIP generates a random IP within the discovered subnet.
func RandomIP() string {
	_, ipNet, err := net.ParseCIDR(subnetPrefix)
	if err != nil {
		// Fallback to a reasonable default
		return fmt.Sprintf("192.0.2.%d", 10+rand.Intn(240))
	}

	// Get the network address as uint32
	networkIP := binary.BigEndian.Uint32(ipNet.IP.To4())

	// Calculate the number of host bits
	ones, bits := ipNet.Mask.Size()
	hostBits := bits - ones

	// Calculate max hosts (excluding network and broadcast)
	maxHosts := (1 << hostBits) - 2
	if maxHosts < 1 {
		maxHosts = 1
	}

	// Generate a random host number (1 to maxHosts, avoiding 0 and broadcast)
	hostNum := uint32(1 + rand.Intn(maxHosts))

	// Combine network and host
	ip := networkIP + hostNum

	// Convert back to IP
	result := make(net.IP, 4)
	binary.BigEndian.PutUint32(result, ip)

	return result.String()
}
