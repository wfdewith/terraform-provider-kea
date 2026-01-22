# Basic DHCP reservation using MAC address
resource "kea_dhcp4_reservation" "server01" {
  subnet_id  = 1
  hw_address = "aa:bb:cc:dd:ee:01"
  ip_address = "192.0.2.10"
}

# Reservation using client ID
resource "kea_dhcp4_reservation" "server02" {
  subnet_id  = 1
  client_id  = "01:aa:bb:cc:dd:ee:02"
  ip_address = "192.0.2.11"
}

# Global reservation (not bound to a specific subnet)
resource "kea_dhcp4_reservation" "global_device" {
  subnet_id  = 0
  hw_address = "aa:bb:cc:dd:ee:ff"
  ip_address = "192.0.2.100"
}
