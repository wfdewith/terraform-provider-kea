package dhcp4_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wfdewith/terraform-provider-kea/internal/acctest"
)

func TestAccReservationDataSource_basic(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	resourceName := "data.kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationDataSourceConfig_basic(acctest.SubnetID(), mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccReservationDataSource_withHostname(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := acctest.RandomIP()
	hostname := fmt.Sprintf("test-host-%s", mac[12:])
	resourceName := "data.kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationDataSourceConfig_withHostname(acctest.SubnetID(), mac, ip, hostname),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
				),
			},
		},
	})
}

func TestAccReservationDataSource_global(t *testing.T) {
	acctest.PreCheck(t)

	mac := acctest.RandomMAC()
	ip := "10.99.99.98"
	resourceName := "data.kea_dhcp4_reservation.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationDataSourceConfig_global(mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subnet_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
				),
			},
		},
	})
}

func testAccReservationDataSourceConfig_basic(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q
}

data "kea_dhcp4_reservation" "test" {
  subnet_id  = kea_dhcp4_reservation.test.subnet_id
  hw_address = kea_dhcp4_reservation.test.hw_address
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationDataSourceConfig_withHostname(subnetID uint32, mac, ip, hostname string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q
  hostname   = %q
}

data "kea_dhcp4_reservation" "test" {
  subnet_id  = kea_dhcp4_reservation.test.subnet_id
  hw_address = kea_dhcp4_reservation.test.hw_address
}
`, acctest.ProviderConfig(), subnetID, mac, ip, hostname)
}

func testAccReservationDataSourceConfig_global(mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = 0
  hw_address = %q
  ip_address = %q
}

data "kea_dhcp4_reservation" "test" {
  subnet_id  = kea_dhcp4_reservation.test.subnet_id
  hw_address = kea_dhcp4_reservation.test.hw_address
}
`, acctest.ProviderConfig(), mac, ip)
}
