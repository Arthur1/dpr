variable "package_store_s3_bucket_name" {
  description = "name of S3 bucket name for package store"
  type        = string
}

variable "tag_db_dynamodb_table_name" {
  description = "name of DynamoDB table name for tag database"
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
