variable "environment" {
  type    = string
  default = "dev"
}

variable "db_username" {
  type      = string
  sensitive = true
}