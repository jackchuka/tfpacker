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