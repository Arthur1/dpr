variable "packages_store_s3_bucket_name" {
  description = "name of S3 bucket name for packages store"
  type        = string
}

variable "tags_db_dynamodb_table_name" {
  description = "name of DynamoDB table name for tags database"
  type        = string
}

variable "tags" {
  type = map(string)
	default = {}
}
