<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.27.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 5.27.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_dynamodb_table.tag_db](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table) | resource |
| [aws_s3_bucket.package_store](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [aws_s3_bucket_ownership_controls.package_store](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_ownership_controls) | resource |
| [aws_s3_bucket_public_access_block.package_store](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_public_access_block) | resource |
| [aws_s3_bucket_server_side_encryption_configuration.package_store](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_server_side_encryption_configuration) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_package_store_s3_bucket_name"></a> [package\_store\_s3\_bucket\_name](#input\_package\_store\_s3\_bucket\_name) | name of S3 bucket name for package store | `string` | n/a | yes |
| <a name="input_tag_db_dynamodb_table_name"></a> [tag\_db\_dynamodb\_table\_name](#input\_tag\_db\_dynamodb\_table\_name) | name of DynamoDB table name for tag database | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | n/a | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_packages_store_s3_bucket_id"></a> [packages\_store\_s3\_bucket\_id](#output\_packages\_store\_s3\_bucket\_id) | ID of S3 bucket for packages store |
| <a name="output_tags_db_dynamodb_table_id"></a> [tags\_db\_dynamodb\_table\_id](#output\_tags\_db\_dynamodb\_table\_id) | ID of DynamoDB table for tags database |
<!-- END_TF_DOCS -->