// Resources with db_ prefix
resource "aws_rds_instance" "db_main" {
  name     = "main-database"
  instance = "db.t3.micro"
  engine   = "postgres"
}

resource "aws_dynamodb_table" "db_users" {
  name         = "users-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
}

// Resources with _db suffix
resource "aws_security_group" "postgres_db" {
  name        = "allow_postgres"
  description = "Allow PostgreSQL inbound traffic"
}

resource "aws_iam_role" "rds_db" {
  name = "rds-monitoring-role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "monitoring.rds.amazonaws.com"
        }
      },
    ]
  })
}

// Resources with api_ prefix
resource "aws_api_gateway_rest_api" "api_main" {
  name        = "main-api"
  description = "Main API Gateway"
}

resource "aws_lambda_function" "api_handler" {
  function_name = "api-handler"
  handler       = "index.handler"
  runtime       = "nodejs14.x"
  role          = aws_iam_role.lambda_role.arn
  
  filename      = "dummy.zip"  // For testing purposes only
}

// Lambda execution role
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

// Variables with db_ prefix
variable "db_username" {
  type      = string
  sensitive = true
}

variable "db_password" {
  type      = string
  sensitive = true
}

// Variable with _db suffix
variable "postgres_db" {
  type    = string
  default = "main"
}

// Regular resources (no special naming)
resource "aws_s3_bucket" "logs" {
  bucket = "my-logs-bucket"
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name = "lambda-logs"
}

// Data sources
data "aws_region" "current" {}

data "aws_availability_zones" "available" {}

// Modules
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  name   = "my-vpc"
}

// Outputs
output "db_endpoint" {
  value = aws_rds_instance.db_main.endpoint
}

output "api_url" {
  value = aws_api_gateway_rest_api.api_main.id
}

// Locals
locals {
  common_tags = {
    Environment = "dev"
    Project     = "test"
  }
}

// Provider
provider "aws" {
  region = "us-west-2"
}
