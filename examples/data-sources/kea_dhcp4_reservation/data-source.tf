# Look up an existing reservation by MAC address
data "kea_dhcp4_reservation" "existing_server" {
  subnet_id  = 1
  hw_address = "aa:bb:cc:dd:ee:01"
}

# Use the data from the existing reservation
output "server_ip" {
  value = data.kea_dhcp4_reservation.existing_server.ip_address
}

output "server_hostname" {
  value = data.kea_dhcp4_reservation.existing_server.hostname
}

output "server_options" {
  value = data.kea_dhcp4_reservation.existing_server.option_data
}

# Look up a global reservation
data "kea_dhcp4_reservation" "global_lookup" {
  subnet_id  = 0
  hw_address = "aa:bb:cc:dd:ee:ff"
}

# Look up using client ID
data "kea_dhcp4_reservation" "by_client_id" {
  subnet_id = 2
  client_id = "01:aa:bb:cc:dd:ee:02"
}

# Look up a reservation by IP address
data "kea_dhcp4_reservation" "by_ip" {
  subnet_id  = 1
  ip_address = "192.0.2.100"
}

# Use the user context from the reservation
data "kea_dhcp4_reservation" "with_metadata" {
  subnet_id  = 3
  hw_address = "aa:bb:cc:dd:ee:50"
}

# Parse and use the user context JSON
locals {
  device_metadata = jsondecode(data.kea_dhcp4_reservation.with_metadata.user_context)
}

output "device_owner" {
  value = local.device_metadata.owner
}

output "device_tags" {
  value = local.device_metadata.tags
}

output "device_location" {
  value = local.device_metadata.metadata.location
}
