# Configure via HTTP endpoint with basic authentication
provider "kea" {
  dhcp4 = {
    address       = "http://192.0.2.1:8000"
    http_username = "admin"
    http_password = "secret"
  }
}

# Configure via UNIX socket
provider "kea" {
  dhcp4 = {
    address = "unix:///run/kea/kea-dhcp4.sock"
  }
}

# Configure via environment variables
# export KEA_DHCP4_ADDRESS="http://192.0.2.1:8000"
# export KEA_DHCP4_HTTP_USERNAME="admin"
# export KEA_DHCP4_HTTP_PASSWORD="secret"
provider "kea" {
}
