provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_activation" "test" {
  property_id = "prp_test"
  version = 1
}