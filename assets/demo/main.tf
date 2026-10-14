variable "region" {
  default = "us-east-1"
}

resource "aws_s3_bucket" "assets" {
  bucket = "example-app-assets"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  cidr   = "10.0.0.0/16"
}

# Primary database for the app
resource "aws_db_instance" "app_db" {
  engine         = "postgres"
  instance_class = "db.t4g.micro" # smallest Graviton class
}

provider "aws" {
  region = var.region
}

output "bucket_arn" {
  value = aws_s3_bucket.assets.arn
}

resource "aws_iam_role" "app_role" {
  name = "example-app"
}

# Set via TF_VAR_db_password
variable "db_password" {
  sensitive = true
}

resource "aws_s3_bucket" "logs" {
  bucket = "example-app-logs"
}

data "aws_caller_identity" "current" {}

resource "aws_iam_policy" "app_policy" {
  name = "example-app"
}

output "db_endpoint" {
  value = aws_db_instance.app_db.endpoint
}

locals {
  app_name = "example-app"
}
