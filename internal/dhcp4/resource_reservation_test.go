package dhcp4_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/wfdewith/terraform-provider-kea/internal/acctest"
	"github.com/wfdewith/terraform-provider-kea/kea"
)

func TestAccReservation_basic(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_basic(acctest.SubnetID(), mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccReservation_withHostname(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	hostname := fmt.Sprintf("test-host-%d", rand.Intn(10000))
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withHostname(acctest.SubnetID(), mac, ip, hostname),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
				),
			},
		},
	})
}

func TestAccReservation_withOptionData(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withOptionData(acctest.SubnetID(), mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "option_data.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "option_data.0.name", "domain-name-servers"),
				),
			},
		},
	})
}

func TestAccReservation_disappears(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_basic(acctest.SubnetID(), mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					testAccCheckReservationDisappears(resourceName, mac),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccReservation_global(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := "10.99.99.99" // Use a fixed IP outside normal subnets for global reservation
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_global(mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccReservation_update(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	hostname1 := fmt.Sprintf("test-host-%d", rand.Intn(10000))
	hostname2 := fmt.Sprintf("test-host-%d", rand.Intn(10000))
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withHostname(acctest.SubnetID(), mac, ip, hostname1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname1),
				),
			},
			{
				Config: testAccReservationConfig_withHostname(acctest.SubnetID(), mac, ip, hostname2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname2),
				),
			},
		},
	})
}

func TestAccReservation_reorderSetsNoUpdate(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withClientClassesAndOptions(acctest.SubnetID(), mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "client_classes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "option_data.#", "3"),
				),
			},
			{
				Config:   testAccReservationConfig_withClientClassesAndOptionsReordered(acctest.SubnetID(), mac, ip),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckReservationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		return nil
	}
}

func testAccCheckReservationDestroy(s *terraform.State) error {
	address := os.Getenv("KEA_DHCP4_ADDRESS")
	transport := &kea.HTTPTransport{Endpoint: address}
	client := kea.NewDHCP4Client(transport)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kea_dhcp4_reservation" {
			continue
		}

		hwAddress := rs.Primary.Attributes["hw_address"]
		subnetIDStr := rs.Primary.Attributes["subnet_id"]
		subnetID := uint32(0)
		if subnetIDStr != "" {
			var parsed int
			fmt.Sscanf(subnetIDStr, "%d", &parsed)
			subnetID = uint32(parsed)
		}

		reservation, err := client.GetReservation(context.Background(), kea.QueryReservationByIdentifier(subnetID, "hw-address", hwAddress))
		if err != nil {
			return err
		}

		if reservation != nil {
			return fmt.Errorf("reservation %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckReservationDisappears(resourceName string, mac string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		address := os.Getenv("KEA_DHCP4_ADDRESS")
		transport := &kea.HTTPTransport{Endpoint: address}
		client := kea.NewDHCP4Client(transport)

		return client.DeleteReservation(context.Background(), kea.QueryReservationByIdentifier(acctest.SubnetID(), "hw-address", mac))
	}
}

func testAccReservationConfig_basic(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationConfig_withHostname(subnetID uint32, mac, ip, hostname string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q
  hostname   = %q
}
`, acctest.ProviderConfig(), subnetID, mac, ip, hostname)
}

func testAccReservationConfig_withOptionData(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q

  option_data {
    name       = "domain-name-servers"
    data       = "8.8.8.8,8.8.4.4"
    csv_format = true
  }
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationConfig_global(mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = 0
  hw_address = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), mac, ip)
}

func testAccReservationConfig_withClientClassesAndOptions(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id      = %d
  hw_address     = %q
  ip_address     = %q
  client_classes = ["web-servers", "production", "monitoring"]

  option_data {
    name = "domain-name-servers"
    data = "8.8.8.8,8.8.4.4"
  }

  option_data {
    name = "domain-name"
    data = "example.com"
  }

  option_data {
    name = "routers"
    data = "192.0.2.1"
  }
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationConfig_withClientClassesAndOptionsReordered(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id      = %d
  hw_address     = %q
  ip_address     = %q
  client_classes = ["monitoring", "production", "web-servers"]

  option_data {
    name = "routers"
    data = "192.0.2.1"
  }

  option_data {
    name = "domain-name-servers"
    data = "8.8.8.8,8.8.4.4"
  }

  option_data {
    name = "domain-name"
    data = "example.com"
  }
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func TestAccReservation_withUserContext(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	resourceName := "kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withUserContext(acctest.SubnetID(), mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReservationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "user_context"),
				),
			},
		},
	})
}

func testAccReservationConfig_withUserContext(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id    = %d
  hw_address   = %q
  ip_address   = %q
  user_context = jsonencode({
    department = "engineering"
    owner      = "network-team"
    tags       = ["production", "critical"]
  })
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}
