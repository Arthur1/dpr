resource "aws_cloudwatch_event_rule" "schedule_rule" {
  name                = var.rule_name
  description         = "created by mackerel-monitoring-modules"
  state               = var.rule_is_enabled ? "ENABLED" : "DISABLED"
  schedule_expression = var.rule_schedule_expression
  tags                = var.tags
}

resource "aws_lambda_permission" "this" {
  action        = "lambda:InvokeFunction"
  function_name = var.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.schedule_rule.arn
}
