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

variable "db_username" {
  type      = string
  sensitive = true
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "postgres_db" {
  type    = string
  default = "main"
}

output "db_endpoint" {
  value = aws_rds_instance.db_main.endpoint
}