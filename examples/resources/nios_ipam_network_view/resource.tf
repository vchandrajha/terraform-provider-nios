terraform {
  required_providers {
    nios = {
      source  = "infobloxopen/nios"
      version = "1.1.0"
    }
  }
}

provider "nios" {
  nios_host_url = "https://172.28.83.72"
  nios_username = "admin"
  nios_password = "Infoblox@123"
}

// Create Network View with Basic Fields
resource "nios_ipam_network_view" "create_network_view" {
  name = "example_network_view"
}

// Create Network View with Additional Fields
resource "nios_ipam_network_view" "create_network_view_with_additional_fields" {
  name    = "example-network-view3"
  comment = "Example Network View with Additional Fields test"

  remote_reverse_zones = [
    {
      fqdn           = "0.168.192.in-addr.arpa"
      key_type       = "NONE"
      server_address = "192.168.12.13"
    },
    {
      fqdn           = "1.168.192.in-addr.arpa"
      key_type       = "TSIG"
      server_address = "192.168.12.13"
      tsig_key_name  = "aeiou"
      tsig_key_alg   = "HMAC-SHA256"
      tsig_key       = "dGhpc2lzdGVzdHRzaWdrZXk="
    }
  ]

  remote_forward_zones = [
    {
      fqdn           = "fwdzone1.com"
      key_type       = "NONE"
      server_address = "192.168.12.13"
    },
    {
      fqdn           = "fwdzone2.com"
      key_type       = "TSIG"
      server_address = "192.168.12.13"
      tsig_key_name  = "aeiou"
      tsig_key_alg   = "HMAC-SHA256"
      tsig_key       = "dGhpc2lzdGVzdHRzaWdrZXk="
    }
  ]
  mgm_private = true

  extattrs = {
    Site = "location-2"
  }
}
