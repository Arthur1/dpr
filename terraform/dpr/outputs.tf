output "packages_store_s3_bucket_id" {
  description = "ID of S3 bucket for packages store"
  value       = aws_s3_bucket.package_store.id
}

output "tags_db_dynamodb_table_id" {
  description = "ID of DynamoDB table for tags database"
  value       = aws_dynamodb_table.tag_db.id
}
