resource "aws_s3_bucket" "logs" {
  bucket = "my-logs-bucket"
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name = "lambda-logs"
}