# variable "vpc_id" {}

data "aws_vpc" "selected" {
	id                 = var.vpc_id
}

variable "regional_cidr_blocks" {
	description = "A simple map of subnets used by region"
	type        = map
	default     = {
		"us-west-2a-private" = "10.172.64.0/23",
		"us-west-2b-private" = "10.172.66.0/23",
		"us-west-2a-public"  = "10.172.68.0/26",
		"us-west-2b-public"  = "10.172.68.64/26",
		"us-east-1a-private" = "10.172.0.0/23",
		"us-east-1b-private" = "10.172.2.0/23",
		"us-east-1a-public"  = "10.172.4.0/26",
		"us-east-1b-public"  = "10.172.4.64/26"
	}
}

resource "aws_subnet" "aza-private" {
	depends_on         = [data.aws_vpc.selected]
	vpc_id             = data.aws_vpc.id
	availability_zone  = format("%sa", var.region)
	cidr_block         = var.regional_cidr_blocks[
		format("%sa-private", var.region)
	]
	tags               = merge(
		var.extra_tags,
		{
			Name = "aza-private"
			tier = "private"
		}
	)
}

# resource "aws_subnet" "aza-public" {
#   depends_on         = [data.aws_vpc.selected]
#   vpc_id             = data.aws_vpc.selected.id
#   availability_zone  = format("%sa", var.region)
#   cidr_block         = var.regional_cidr_blocks[
#     format("%sa-public", var.region)
#   ]
#   tags               = merge(
#     var.extra_tags,
#     {
#       Name = "aza-public"
#       tier = "public"
#     }
#   )
# }
# 
# resource "aws_subnet" "azb-private" {
#   depends_on         = [data.aws_vpc.selected]
#   vpc_id             = data.aws_vpc.selected.id
#   availability_zone  = format("%sb", var.region)
#   cidr_block         = var.regional_cidr_blocks[
#     format("%sb-private", var.region)
#   ]
#   tags               = merge(
#     var.extra_tags,
#     {
#       Name = "azb-private"
#       tier = "private"
#     }
#   )
# }
# 
# resource "aws_subnet" "azb-public" {
#   depends_on         = [data.aws_vpc.selected]
#   vpc_id             = data.aws_vpc.selected.id
#   availability_zone  = format("%sb", var.region)
#   cidr_block         = var.regional_cidr_blocks[
#     format("%sb-public", var.region)
#   ]
#   tags               = merge(
#     var.extra_tags,
#     {
#       Name = "azb-public"
#       tier = "public"
#     }
#   )
# }
