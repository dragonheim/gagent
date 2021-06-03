variable "extra_tags" {
	description        = "Tags required on all resources"
	type               = map
	default            = {
		"org"            = "dragonheim"
		"service"        = "gagent"
		"maintained_by"  = "jwells@dragonheim.net"
	}
}
