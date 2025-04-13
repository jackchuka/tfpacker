resource "aws_s3_bucket" "bucket1" {
  bucket = "my-test-bucket-1"

  tags = {
    Name        = "Test Bucket 1"
    Environment = "Dev"
  }
}