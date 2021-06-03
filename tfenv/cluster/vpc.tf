variable "region" {}
variable "provider_alias" {}

variable "regional_vpc_cidr" {
	description = "A simple map of VPC subnets used by region"
	type        = map
	default     = {
		"us-west-2" = "10.172.64.0/19",
		"us-east-1" = "10.172.0.0/19",
	}
}

resource "aws_vpc" "gagent" {
	instance_tenancy   = "default"
	enable_dns_support = true
	cidr_block         = var.regional_vpc_cidr[var.region]
	tags = merge(
		var.extra_tags,
		{
			Name = "gagent"
		}
	)
}
