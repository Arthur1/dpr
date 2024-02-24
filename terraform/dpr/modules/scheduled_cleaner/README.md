<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.27.0 |
| <a name="requirement_github"></a> [github](#requirement\_github) | ~> 6.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.27.0 |
| <a name="provider_github"></a> [github](#provider\_github) | ~> 6.0 |
| <a name="provider_terraform"></a> [terraform](#provider\_terraform) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_event_rule.schedule_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_rule) | resource |
| [aws_iam_role.scheduled_cleaner](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.scheduled_cleaner](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy_attachment.lambda_basic_execution](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_lambda_function.cleaner_cron](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function) | resource |
| [aws_lambda_permission.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_permission) | resource |
| [terraform_data.dpr_cleaner_eventbridge_lambda](https://registry.terraform.io/providers/hashicorp/terraform/latest/docs/resources/data) | resource |
| [aws_iam_policy.lambda_basic_execution](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy) | data source |
| [aws_iam_policy_document.assume_role_lambda](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.scheduled_cleaner](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [github_release.dpr](https://registry.terraform.io/providers/integrations/github/latest/docs/data-sources/release) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_function_memory_size"></a> [function\_memory\_size](#input\_function\_memory\_size) | n/a | `number` | `128` | no |
| <a name="input_function_name"></a> [function\_name](#input\_function\_name) | n/a | `string` | n/a | yes |
| <a name="input_function_timeout"></a> [function\_timeout](#input\_function\_timeout) | n/a | `number` | `7` | no |
| <a name="input_package_store_s3_bucket_id"></a> [package\_store\_s3\_bucket\_id](#input\_package\_store\_s3\_bucket\_id) | id pf S3 bucketfor package store | `string` | n/a | yes |
| <a name="input_rule_is_enabled"></a> [rule\_is\_enabled](#input\_rule\_is\_enabled) | n/a | `bool` | `true` | no |
| <a name="input_rule_name"></a> [rule\_name](#input\_rule\_name) | n/a | `string` | `"dpr-cleaner-schedule"` | no |
| <a name="input_rule_schedule_expression"></a> [rule\_schedule\_expression](#input\_rule\_schedule\_expression) | n/a | `string` | `"rate(24 hour)"` | no |
| <a name="input_tag_db_dynamodb_table_id"></a> [tag\_db\_dynamodb\_table\_id](#input\_tag\_db\_dynamodb\_table\_id) | id of DynamoDB table for tag database | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | n/a | `map(string)` | `{}` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->