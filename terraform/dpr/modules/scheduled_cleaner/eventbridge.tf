resource "aws_cloudwatch_event_rule" "schedule_rule" {
  name                = var.rule_name
  description         = "created by dpr-scheduled-cleaner"
  state               = var.rule_is_enabled ? "ENABLED" : "DISABLED"
  schedule_expression = var.rule_schedule_expression
  tags                = var.tags
}

resource "aws_cloudwatch_event_target" "schedule_rule" {
  rule = aws_cloudwatch_event_rule.schedule_rule.name
  arn  = aws_lambda_function.scheduled_cleaner.arn

  input = jsonencode({
    "package-store" = {
      "s3-bucket-name" = var.package_store_s3_bucket_name
    }
    "tag-db" = {
      "dynamodb-table-name" = var.tag_db_dynamodb_table_name
    }
    "lifecycle-policy" = yamldecode(file(var.lifecycle_policy_file_path))
  })
}
