module "scheduled_cleaner" {
  source                       = "./modules/scheduled_cleaner"
  package_store_s3_bucket_name = aws_s3_bucket.package_store.bucket
  tag_db_dynamodb_table_name   = aws_dynamodb_table.tag_db.name
  lifecycle_policy_file_path   = var.lifecycle_policy_file_path
  function_name                = "dpr-scheduled-cleaner"
  rule_schedule_expression     = var.scheduled_cleaner_schedule_expression
  tags                         = var.tags
}
