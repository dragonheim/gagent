# main.tf
module "us-east-1" {
	source    = "./cluster"
	region    = "us-east-1"
	provider_alias = us-west-2
	providers = {
		aws = "aws.us-east-1"
	}
}

module "us-west-2" {
	source    = "./cluster"
	region    = "us-west-2"
	provider_alias = us-west-2
	providers = {
		aws = "aws.us-west-2"
	}
}
