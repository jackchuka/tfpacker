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