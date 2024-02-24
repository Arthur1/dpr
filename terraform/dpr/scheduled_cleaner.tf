module "scheduled_cleaner" {
  source                     = "./modules/scheduled_cleaner"
  package_store_s3_bucket_id = aws_s3_bucket.package_store.id
  tag_db_dynamodb_table_id   = aws_dynamodb_table.tag_db.id
  function_name              = "dpr-scheduled-cleaner"
  rule_schedule_expression   = "rate(5 min)"
  tags                       = var.tags
}
