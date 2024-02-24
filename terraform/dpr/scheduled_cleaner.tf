module "scheduled_cleaner" {
  source                      = "./modules/scheduled_cleaner"
  package_store_s3_bucket_arn = aws_s3_bucket.package_store.arn
  tag_db_dynamodb_table_arn   = aws_dynamodb_table.tag_db.arn
  function_name               = "dpr-scheduled-cleaner"
  rule_schedule_expression    = "rate(5 minutes)"
  tags                        = var.tags
}
