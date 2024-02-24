output "package_store_s3_bucket_id" {
  description = "ID of S3 bucket for packages store"
  value       = aws_s3_bucket.package_store.id
}

output "package_store_s3_bucket_arn" {
  description = "arn of S3 bucket for packages store"
  value       = aws_s3_bucket.package_store.arn
}

output "tag_db_dynamodb_table_id" {
  description = "ID of DynamoDB table for tags database"
  value       = aws_dynamodb_table.tag_db.id
}

output "tag_db_dynamodb_table_arn" {
  description = "arn of DynamoDB table for tags database"
  value       = aws_dynamodb_table.tag_db.arn
}
