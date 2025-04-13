module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  name   = "shared-vpc"
}