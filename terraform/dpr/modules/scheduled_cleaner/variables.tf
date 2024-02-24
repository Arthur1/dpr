variable "package_store_s3_bucket_id" {
  description = "id pf S3 bucketfor package store"
  type        = string
}

variable "tag_db_dynamodb_table_id" {
  description = "id of DynamoDB table for tag database"
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
  default = "rate(24 hour)"
}

variable "tags" {
  type    = map(string)
  default = {}
}

