variable "package_store_s3_bucket_name" {
  description = "name of S3 bucket for package store"
  type        = string
}

variable "tag_db_dynamodb_table_name" {
  description = "name of DynamoDB table for tag database"
  type        = string
}

variable "lifecycle_policy_file_path" {
  type = string
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
  default = "rate(24 hours)"
}

variable "tags" {
  type    = map(string)
  default = {}
}

