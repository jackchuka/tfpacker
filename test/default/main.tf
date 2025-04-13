// this is comment
resource "aws_s3_bucket" "bucket1" {
  bucket = "my-test-bucket-1"

  tags = {
    Name        = "Test Bucket 1"
    Environment = "Dev"
  }
}

// comment
resource "aws_iam_role" "role1" {
  name = "test-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}

data "aws_region" "current" {}

variable "environment" {
  type    = string
  default = "dev"
}

variable "db_username" {
  type      = string
  sensitive = true
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "my-vpc"
  cidr = "10.0.0.0/16"
}

output "bucket_name" {
  value = aws_s3_bucket.bucket1.bucket
}

locals {
  common_tags = {
    Project     = "Test"
    Environment = var.environment
  }
}

provider "aws" {
  region = "us-west-2"
}
