// Production resources
resource "aws_s3_bucket" "prod_logs" {
  bucket = "prod-logs-bucket"
}

resource "aws_dynamodb_table" "prod_users" {
  name         = "prod-users-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
}

resource "aws_lambda_function" "prod_api" {
  function_name = "prod-api-handler"
  handler       = "index.handler"
  runtime       = "nodejs14.x"
  role          = aws_iam_role.lambda_role.arn
  
  filename      = "dummy.zip"  // For testing purposes only
}

variable "prod_region" {
  type    = string
  default = "us-west-2"
}

output "prod_api_url" {
  value = aws_lambda_function.prod_api.function_name
}

// Staging resources
resource "aws_s3_bucket" "staging_logs" {
  bucket = "staging-logs-bucket"
}

resource "aws_dynamodb_table" "staging_users" {
  name         = "staging-users-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
}

resource "aws_lambda_function" "staging_api" {
  function_name = "staging-api-handler"
  handler       = "index.handler"
  runtime       = "nodejs14.x"
  role          = aws_iam_role.lambda_role.arn
  
  filename      = "dummy.zip"  // For testing purposes only
}

variable "staging_region" {
  type    = string
  default = "us-east-1"
}

output "staging_api_url" {
  value = aws_lambda_function.staging_api.function_name
}

// Development resources
resource "aws_s3_bucket" "dev_logs" {
  bucket = "dev-logs-bucket"
}

resource "aws_dynamodb_table" "dev_users" {
  name         = "dev-users-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
}

resource "aws_lambda_function" "dev_api" {
  function_name = "dev-api-handler"
  handler       = "index.handler"
  runtime       = "nodejs14.x"
  role          = aws_iam_role.lambda_role.arn
  
  filename      = "dummy.zip"  // For testing purposes only
}

variable "dev_region" {
  type    = string
  default = "us-east-2"
}

output "dev_api_url" {
  value = aws_lambda_function.dev_api.function_name
}

// Shared resources
resource "aws_iam_role" "lambda_role" {
  name = "lambda-execution-role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  name   = "shared-vpc"
}

data "aws_region" "current" {}

locals {
  environments = ["prod", "staging", "dev"]
}

provider "aws" {
  region = "us-west-2"
}
