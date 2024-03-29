
resource "aws_lambda_function" "scheduled_cleaner" {
  function_name = var.function_name
  runtime       = "provided.al2023"
  handler       = "main"
  role          = aws_iam_role.scheduled_cleaner.arn
  architectures = ["arm64"]
  memory_size   = var.function_memory_size
  timeout       = var.function_timeout
  package_type  = "Zip"
  filename      = "./dpr-cleaner-eventbridge-lambda.zip"
  tags          = var.tags
  depends_on = [
    terraform_data.dpr_cleaner_eventbridge_lambda
  ]
}

resource "aws_lambda_permission" "schedule_rule" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.scheduled_cleaner.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.schedule_rule.arn
}
