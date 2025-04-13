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