resource "aws_dynamodb_table" "tags_db" {
  name           = var.tags_db_dynamodb_table_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 3
  write_capacity = 3

  hash_key = "tag"
  attribute {
    name = "tag"
    type = "S"
  }

  tags = var.tags
}
