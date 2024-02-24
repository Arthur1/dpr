variable "package_store_s3_bucket_arn" {
  description = "arn of S3 bucket for package store"
  type        = string
}

variable "tag_db_dynamodb_table_arn" {
  description = "arn of DynamoDB table for tag database"
  type        = string
}

variable "function_name" {
  type = string
}

variable "function_timeout" {
  type    = number
  default = 7
}

variable "function_memory_size" {
  type    = number
  default = 128
}

variable "rule_name" {
  type    = string
  default = "dpr-cleaner-schedule"
}

variable "rule_is_enabled" {
  type    = bool
  default = true
}

variable "rule_schedule_expression" {
  type    = string
  default = "rate(5 minutes)"
}

variable "tags" {
  type    = map(string)
  default = {}
}

