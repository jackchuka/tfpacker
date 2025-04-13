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