
resource "aws_lambda_function" "cleaner_cron" {
  function_name    = var.function_name
  runtime          = "provided.al2023"
  handler          = "main"
  role             = aws_iam_role.cleaner_cron.arn
  architectures    = ["arm64"]
  memory_size      = var.function_memory_size
  timeout          = var.function_timeout
  package_type     = "Zip"
  filename         = "./dpr-cleaner-eventbridge-lambda.zip"
  source_code_hash = filebase64sha256("./dpr-cleaner-eventbridge-lambda.zip")
  tags             = var.tags
}
